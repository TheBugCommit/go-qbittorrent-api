package qbittorrent

import (
	"context"
	"net/http"
)

// Version returns the qBittorrent version, for example "v5.2.4".
func (c *Client) Version(ctx context.Context) (string, error) {
	return c.requestText(ctx, http.MethodGet, buildPath("/app/version"), nil, defaultRetries)
}

// WebAPIVersion returns the qBittorrent Web API version, for example "2.15.1".
func (c *Client) WebAPIVersion(ctx context.Context) (string, error) {
	return c.requestText(ctx, http.MethodGet, buildPath("/app/webapiVersion"), nil, defaultRetries)
}
