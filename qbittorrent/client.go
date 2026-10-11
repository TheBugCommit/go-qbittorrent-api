package qbittorrent

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/publicsuffix"
)

// Client is a qBittorrent Web API client. Create one with NewClient. A Client
// is intended to be safe for concurrent use by multiple goroutines.
//
// When both an API key and a username/password are set, the API key takes
// precedence and requests are authenticated with an "Authorization: Bearer"
// header.
type Client struct {
	baseURL    string
	httpClient *http.Client
	username   string
	password   string
	apiKey     string
	mu         sync.Mutex
}

// Option configures a Client. Options are applied by NewClient in the order
// they are provided.
type Option func(*Client) error

// defaultTimeout is the HTTP timeout applied to requests when no other
// timeout is configured.
const defaultTimeout = 10 * time.Second

// NewClient creates a Client for the given base URL, e.g.
// "http://localhost:8080". The URL is validated and normalized: a trailing
// slash, if any, is removed.
//
// Options are applied in order; if any of them fails, NewClient returns the
// error and no client is returned.
func NewClient(baseURL string, opts ...Option) (*Client, error) {
	baseURL, err := normalizeBaseURL(baseURL)
	if err != nil {
		return nil, err
	}

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, fmt.Errorf("qbittorrent: create cookie jar: %w", err)
	}

	client := &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
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
	c.httpClient.CloseIdleConnections()
	return nil
}

// WithAPIKey configures the client to authenticate every request with the
// given API key using an "Authorization: Bearer" header. The API key takes
// precedence over credentials set with WithUsernamePassword.
func WithAPIKey(apiKey string) Option {
	return func(c *Client) error {
		c.apiKey = apiKey
		return nil
	}
}

// WithUsernamePassword configures the client to authenticate with a username
// and password through Login. It has no effect when an API key is set.
func WithUsernamePassword(username, password string) Option {
	return func(c *Client) error {
		c.username = username
		c.password = password
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
		c.httpClient = client
		return nil
	}
}

// WithTimeout overrides the HTTP timeout applied to all requests. When the
// HTTP client is replaced with WithHTTPClient, apply WithTimeout afterwards
// so the timeout is set on the new client.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) error {
		c.httpClient.Timeout = d
		return nil
	}
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
