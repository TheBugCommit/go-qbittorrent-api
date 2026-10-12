package qbittorrent

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

const apiBase = "/api/v2"

// Login authenticates the client with the username and password configured
// with WithLogin. It is a no-op when an API key is used, since API keys
// authenticate every request on their own.
//
// Since qBittorrent 5.0 the Web API is JSON based: a successful login returns
// HTTP 204 with an empty body and stores the session cookie in the client's
// cookie jar. Rejected credentials are reported by the server with HTTP 401.
func (c *Client) Login(ctx context.Context) error {
	if c.auth.apiKey != "" {
		return nil
	}

	err := c.login(ctx)

	// Remember the outcome so that send does not retry a failing login on every
	// unauthorized response, which would get the client IP banned.
	c.authFailed.Store(err != nil)

	return err
}

func (c *Client) login(ctx context.Context) error {
	if c.auth.username == "" {
		return fmt.Errorf("%w: no username configured, use WithLogin", ErrInvalidCredentials)
	}

	params := url.Values{
		"username": {c.auth.username},
		"password": {c.auth.password},
	}

	if _, err := c.requestText(ctx, http.MethodPost, buildPath("/auth/login"), params, 0); err != nil {
		// A 401 during login means the credentials were rejected.
		var apiErr *APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusUnauthorized {
			message := apiErr.Message
			if message == "" {
				message = "invalid username or password"
			}
			return newAPIErrorWithSentinel(http.StatusUnauthorized, ErrInvalidCredentials, "%s", message)
		}
		return err
	}

	return nil
}

// Logout logs the current session out. It is a no-op when an API key is used,
// because API keys cannot access the auth endpoints.
func (c *Client) Logout(ctx context.Context) error {
	if c.auth.apiKey != "" {
		return nil
	}

	return c.requestJSON(ctx, http.MethodPost, buildPath("/auth/logout"), nil, nil, 0)
}

func buildPath(path string) string {
	return apiBase + path
}
