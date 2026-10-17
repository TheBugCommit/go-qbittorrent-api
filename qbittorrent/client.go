package qbittorrent

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	"golang.org/x/net/publicsuffix"
)

// config holds the connection settings of a Client.
type config struct {
	url     string
	referer string
}

// auth holds the credentials of a Client. When apiKey is set it takes
// precedence over username and password.
type auth struct {
	apiKey   string
	username string
	password string
}

// Client is a qBittorrent Web API client. Create one with New. A Client
// is intended to be safe for concurrent use by multiple goroutines.
//
// When both an API key and a username/password are set, the API key takes
// precedence and requests are authenticated with an "Authorization: Bearer"
// header.
type Client struct {
	config *config
	http   *http.Client
	auth   auth

	// authFailed records whether the last login attempt failed. It prevents
	// send from retrying the login on every unauthorized response, which would
	// make qBittorrent ban the client IP after too many failed attempts.
	authFailed atomic.Bool
}

// Option configures a Client. Options are applied by New in the order
// they are provided.
type Option func(*Client) error

const (
	defaultTimeout  = 30 * time.Second
	defaultRetries  = 1
	userAgent       = "go-qbittorrent-api"
	apiKeyPrefix    = "qbt_"
	apiKeyLength    = 32
	maxErrBody      = 4 << 10  // 4 KiB
	maxResponseBody = 10 << 20 // 10 MiB
)

// New creates a Client for the given base URL, e.g. "http://localhost:8080".
// The URL is validated and normalized: a trailing slash, if any, is removed.
//
// Options are applied in order; if any of them fails, New returns the error
// and no client is returned.
func New(baseURL string, opts ...Option) (*Client, error) {
	baseURL, err := normalizeBaseURL(baseURL)
	if err != nil {
		return nil, err
	}

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, fmt.Errorf("qbittorrent: create cookie jar: %w", err)
	}

	client := &Client{
		config: &config{
			url:     baseURL,
			referer: baseURL,
		},
		http: &http.Client{
			Timeout: defaultTimeout,
			Jar:     jar,
		},
	}

	for _, opt := range opts {
		if err := opt(client); err != nil {
			return nil, err
		}
	}

	return client, nil
}

// Close releases idle connections held by the underlying HTTP client. It is
// safe to call Close multiple times.
func (c *Client) Close() error {
	c.http.CloseIdleConnections()
	return nil
}

// BaseURL returns the normalized base URL the client was created with.
func (c *Client) BaseURL() string {
	return c.config.url
}

// Referer returns the Referer header sent with every request.
func (c *Client) Referer() string {
	return c.config.referer
}

// WithAPIKey configures the client to authenticate every request with the
// given API key using an "Authorization: Bearer" header. The API key takes
// precedence over credentials set with WithLogin.
func WithAPIKey(apiKey string) Option {
	return func(c *Client) error {
		if err := validateAPIKey(apiKey); err != nil {
			return err
		}
		c.auth.apiKey = apiKey
		return nil
	}
}

// WithLogin configures the client to authenticate with a username
// and password through Login. It has no effect when an API key is set.
func WithLogin(username, password string) Option {
	return func(c *Client) error {
		c.auth.username = username
		c.auth.password = password
		return nil
	}
}

// WithHTTPClient replaces the default HTTP client. The caller keeps ownership
// of the provided client. If it has no cookie jar, one is added so that login
// sessions remain valid across requests.
func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) error {
		if client.Jar == nil {
			jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
			if err != nil {
				return fmt.Errorf("qbittorrent: create cookie jar: %w", err)
			}
			client.Jar = jar
		}
		c.http = client
		return nil
	}
}

// WithTimeout overrides the HTTP timeout applied to all requests. When the
// HTTP client is replaced with WithHTTPClient, apply WithTimeout afterwards
// so the timeout is set on the new client.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) error {
		c.http.Timeout = d
		return nil
	}
}

// WithReferer overrides the Referer header sent with every request. qBittorrent
// rejects requests whose Referer host does not match the WebUI host, so this is
// only needed when the client must advertise a different host than the base
// URL. By default the base URL is used.
func WithReferer(referer string) Option {
	return func(c *Client) error {
		c.config.referer = strings.TrimRight(referer, "/")
		return nil
	}
}

func (c *Client) setHeaders(req *http.Request, contentType string) {
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	req.Header.Set("User-Agent", userAgent)
	if c.config.referer != "" {
		req.Header.Set("Referer", c.config.referer)
	}
	if c.auth.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.auth.apiKey)
	}
}

// newRequest builds an HTTP request for the given API URI. GET parameters are
// sent in the query string and POST parameters in a form-encoded body.
func (c *Client) newRequest(ctx context.Context, method, uri string, val url.Values) (*http.Request, error) {
	if val == nil {
		val = url.Values{}
	}

	var body io.Reader
	var contentType string
	if method == http.MethodPost {
		body = strings.NewReader(val.Encode())
		contentType = "application/x-www-form-urlencoded"
	}

	req, err := http.NewRequestWithContext(ctx, method, uri, body)
	if err != nil {
		return nil, err
	}

	if method != http.MethodPost {
		req.URL.RawQuery = val.Encode()
	}

	c.setHeaders(req, contentType)

	return req, nil
}

// newMultipartRequest builds a POST multipart/form-data request with a file.
func (c *Client) newMultipartRequest(ctx context.Context, uri string, fieldName string, filesPaths []string, val url.Values) (*http.Request, error) {
	data, contentType, err := buildMultipart(fieldName, filesPaths, val)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, bytes.NewReader(data.Bytes()))
	if err != nil {
		return nil, err
	}

	c.setHeaders(req, contentType)
	return req, nil
}

// send executes req and returns the raw response body, always bounded by
// maxResponseBody. When the cookie session has expired it transparently logs in
// again once, provided credentials are configured.
func (c *Client) send(ctx context.Context, req *http.Request, retries int) ([]byte, error) {
	response, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = response.Body.Close() }()

	if isUnauthorized(response.StatusCode) && c.canRelogin(retries) {
		_, _ = io.Copy(io.Discard, response.Body)

		if err := c.Login(ctx); err != nil {
			return nil, err
		}
		if err := resetBody(req); err != nil {
			return nil, err
		}

		return c.send(ctx, req, retries-1)
	}

	if !isSuccess(response.StatusCode) {
		return nil, newAPIError(response.StatusCode, "%s", readErrBody(response.Body))
	}

	return io.ReadAll(io.LimitReader(response.Body, maxResponseBody))
}

// canRelogin reports whether the client may transparently log in again after an
// unauthorized response. It never applies to API key authentication, which
// cannot use the auth endpoints, and it requires configured credentials. It also
// stops retrying once a login attempt has failed, so a wrong password does not
// turn into a burst of failed logins that gets the client IP banned.
func (c *Client) canRelogin(retries int) bool {
	return retries > 0 && c.auth.apiKey == "" && c.auth.username != "" && !c.authFailed.Load()
}

// resetBody replaces req.Body with a fresh copy so the request can be sent
// again. It is a no-op for requests without a body.
func resetBody(req *http.Request) error {
	if req.GetBody == nil {
		return nil
	}

	body, err := req.GetBody()
	if err != nil {
		return err
	}
	req.Body = body

	return nil
}

func (c *Client) getRequest(ctx context.Context, uri string, val url.Values, v any) error {
	return c.requestJSON(ctx, http.MethodGet, uri, val, v, defaultRetries)
}

func (c *Client) postRequest(ctx context.Context, uri string, val url.Values, v any) error {
	return c.requestJSON(ctx, http.MethodPost, uri, val, v, defaultRetries)
}

func (c *Client) postFileRequest(ctx context.Context, uri, fieldName string, filesPaths []string, val url.Values) error {
	req, err := c.newMultipartRequest(ctx, c.config.url+uri, fieldName, filesPaths, val)
	if err != nil {
		return err
	}

	_, err = c.send(ctx, req, defaultRetries)
	return err
}

// requestJSON performs a request and decodes a JSON response into v. A nil v or
// an empty body is treated as success.
func (c *Client) requestJSON(ctx context.Context, method, uri string, val url.Values, v any, retries int) error {
	req, err := c.newRequest(ctx, method, c.config.url+uri, val)
	if err != nil {
		return err
	}

	data, err := c.send(ctx, req, retries)
	if err != nil {
		return err
	}

	if v == nil || len(data) == 0 {
		return nil
	}
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("%w: unmarshal JSON: %w", ErrInvalidResponse, err)
	}

	return nil
}

// requestText performs a request whose response body is plain text, such as the
// login and version endpoints, and returns the trimmed body.
func (c *Client) requestText(ctx context.Context, method, uri string, val url.Values, retries int) (string, error) {
	req, err := c.newRequest(ctx, method, c.config.url+uri, val)
	if err != nil {
		return "", err
	}

	data, err := c.send(ctx, req, retries)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(data)), nil
}

func buildMultipart(fieldName string, paths []string, val url.Values) (*bytes.Buffer, string, error) {
	buf := &bytes.Buffer{}
	mpw := multipart.NewWriter(buf)
	defer func() { _ = mpw.Close() }()

	for _, path := range paths {
		if err := attachFile(mpw, fieldName, path); err != nil {
			return nil, "", err
		}
	}

	for key, values := range val {
		for _, value := range values {
			if err := mpw.WriteField(key, value); err != nil {
				return nil, "", err
			}
		}
	}

	if err := mpw.Close(); err != nil {
		return nil, "", err
	}

	return buf, mpw.FormDataContentType(), nil
}

func attachFile(mpw *multipart.Writer, fieldName string, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	fileWriter, err := mpw.CreateFormFile(fieldName, filepath.Base(path))
	if err != nil {
		return err
	}

	if _, err = io.Copy(fileWriter, file); err != nil {
		return err
	}

	return nil
}

func isUnauthorized(status int) bool {
	return status == http.StatusUnauthorized || status == http.StatusForbidden
}

func isSuccess(status int) bool {
	return status >= http.StatusOK && status < http.StatusMultipleChoices
}

func readErrBody(body io.Reader) string {
	buf, err := io.ReadAll(io.LimitReader(body, maxErrBody))
	if err != nil {
		return err.Error()
	}

	return strings.TrimSpace(string(buf))
}

// validateAPIKey checks that apiKey matches the format used by qBittorrent:
// the "qbt_" prefix followed by 28 alphanumeric characters.
func validateAPIKey(apiKey string) error {
	if len(apiKey) != apiKeyLength || !strings.HasPrefix(apiKey, apiKeyPrefix) {
		return fmt.Errorf("%w: expected %q followed by %d alphanumeric characters",
			ErrInvalidAPIKeyFormat, apiKeyPrefix, apiKeyLength-len(apiKeyPrefix))
	}

	for _, r := range apiKey[len(apiKeyPrefix):] {
		if !isAlphanumeric(r) {
			return fmt.Errorf("%w: unexpected character %q", ErrInvalidAPIKeyFormat, r)
		}
	}

	return nil
}

func isAlphanumeric(r rune) bool {
	return r >= '0' && r <= '9' || r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z'
}

// normalizeBaseURL validates the given base URL and returns it without a
// trailing slash. Only http and https URLs with a hostname are accepted.
func normalizeBaseURL(rawURL string) (string, error) {
	u, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return "", fmt.Errorf("qbittorrent: invalid base URL %q: %w", rawURL, err)
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return "", fmt.Errorf("qbittorrent: base URL scheme must be http or https, got %q", u.Scheme)
	}
	if u.Hostname() == "" {
		return "", errors.New("qbittorrent: base URL is missing a hostname")
	}
	if u.RawQuery != "" || u.Fragment != "" {
		return "", errors.New("qbittorrent: base URL must not contain a query or fragment")
	}

	return strings.TrimRight(rawURL, "/"), nil
}
