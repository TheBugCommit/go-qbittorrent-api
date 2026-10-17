package qbittorrent

import (
	"context"

	"github.com/TheBugCommit/go-qbittorrent-api/models"
)

// GetTorrents returns the torrents matching the given filters. A zero-value
// TorrentListRequest returns every torrent.
func (c *Client) GetTorrents(ctx context.Context, req models.TorrentListRequest) ([]models.Torrent, error) {
	var torrents []models.Torrent
	if err := c.getRequest(ctx, buildPath("/torrents/info"), req.ToValues(), &torrents); err != nil {
		return nil, err
	}

	return torrents, nil
}

// GetTorrentProperties returns the detailed properties of a single torrent.
func (c *Client) GetTorrentProperties(ctx context.Context, hash string) (*models.TorrentProperties, error) {
	req := models.TorrentHashRequest{Hash: hash}

	props := &models.TorrentProperties{}
	if err := c.getRequest(ctx, buildPath("/torrents/properties"), req.ToValues(), props); err != nil {
		return nil, err
	}

	return props, nil
}

// AddTorrent adds one or more torrents from URLs o File paths or magnet links.
func (c *Client) AddTorrent(ctx context.Context, req models.AddTorrentRequest) error {
	uri := buildPath("/torrents/add")

	if len(req.Files) > 0 {
		return c.postFileRequest(ctx, uri, "torrents", req.Files, req.ToValues())

	}

	return c.postRequest(ctx, uri, req.ToValues(), nil)
}

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
