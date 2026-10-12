package qbittorrent

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/TheBugCommit/go-qbittorrent-api/models"
)

// GetTransferInfo returns the global transfer information.
func (c *Client) GetTransferInfo(ctx context.Context) (*models.TransferInfo, error) {
	info := &models.TransferInfo{}
	if err := c.getRequest(ctx, buildPath("/transfer/info"), nil, info); err != nil {
		return nil, err
	}

	return info, nil
}

// GetDownloadLimit returns the global download speed limit in bytes per second.
// A limit of 0 means unlimited.
//
// The server answers with a bare number (for example "0"), not a JSON object.
func (c *Client) GetDownloadLimit(ctx context.Context) (int64, error) {
	return c.getTransferLimit(ctx, "/transfer/downloadLimit")
}

// GetUploadLimit returns the global upload speed limit in bytes per second.
// A limit of 0 means unlimited.
//
// The server answers with a bare number (for example "0"), not a JSON object.
func (c *Client) GetUploadLimit(ctx context.Context) (int64, error) {
	return c.getTransferLimit(ctx, "/transfer/uploadLimit")
}

// SetDownloadLimit sets the global download speed limit in bytes per second.
// A limit of 0 means unlimited.
func (c *Client) SetDownloadLimit(ctx context.Context, req models.SetTransferLimitRequest) error {
	return c.setTransferLimit(ctx, "/transfer/setDownloadLimit", req)
}

// SetUploadLimit sets the global upload speed limit in bytes per second.
// A limit of 0 means unlimited.
func (c *Client) SetUploadLimit(ctx context.Context, req models.SetTransferLimitRequest) error {
	return c.setTransferLimit(ctx, "/transfer/setUploadLimit", req)
}

func (c *Client) getTransferLimit(ctx context.Context, path string) (int64, error) {
	body, err := c.requestText(ctx, http.MethodGet, buildPath(path), nil, defaultRetries)
	if err != nil {
		return 0, err
	}

	limit, err := strconv.ParseInt(body, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: parse transfer limit %q: %w", ErrInvalidResponse, body, err)
	}

	return limit, nil
}

func (c *Client) setTransferLimit(ctx context.Context, path string, req models.SetTransferLimitRequest) error {
	return c.postRequest(ctx, buildPath(path), req.ToValues(), nil)
}
