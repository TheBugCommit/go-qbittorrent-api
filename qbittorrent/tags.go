package qbittorrent

import (
	"context"

	"github.com/TheBugCommit/go-qbittorrent-api/models"
)

// GetTags returns all qBittorrent client tags.
//
// The tags live under the torrents scope, so the endpoint is
// /api/v2/torrents/tags and the response is a JSON array of tag names.
func (c *Client) GetTags(ctx context.Context) ([]models.Tag, error) {
	var tags []string
	if err := c.getRequest(ctx, buildPath("/torrents/tags"), nil, &tags); err != nil {
		return nil, err
	}

	return convertToTags(tags), nil
}

// CreateTags creates one or more tags.
func (c *Client) CreateTags(ctx context.Context, req models.CreateTagRequest) error {
	return c.postRequest(ctx, buildPath("/torrents/createTags"), req.ToValues(), nil)
}

// DeleteTags deletes one or more tags.
func (c *Client) DeleteTags(ctx context.Context, req models.DeleteTagRequest) error {
	return c.postRequest(ctx, buildPath("/torrents/deleteTags"), req.ToValues(), nil)
}

func convertToTags(names []string) []models.Tag {
	tags := make([]models.Tag, len(names))

	for i := range names {
		tags[i] = models.Tag{Name: names[i]}
	}

	return tags
}
