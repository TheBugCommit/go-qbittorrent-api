package qbittorrent

import (
	"context"
	"strings"

	"github.com/TheBugCommit/go-qbittorrent-api/models"
)

//*
//TODO: Sync, Delete, Reecheck, Reeanunnce,SetTorrentCategory, SetAutomaticTorrentManagment, AddTorrentTags,RemoveTorrentTags

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

// PauseTorrents pause one or more torrents. Pause all torrents, without passing torrent hashes.
func (c *Client) PauseTorrents(ctx context.Context, req models.TorrentHashesRequest) error {
	val := req.ToValues()
	if strings.TrimSpace(val.Get("hashes")) == "" {
		val.Set("hashes", "all")
	}

	return c.postRequest(ctx, buildPath("/torrents/stop"), val, nil)
}

// ResumeTorrents resume one or more torrents. Resume all torrents, without passing torrent hashes.
func (c *Client) ResumeTorrents(ctx context.Context, req models.TorrentHashesRequest) error {
	val := req.ToValues()
	if strings.TrimSpace(val.Get("hashes")) == "" {
		val.Set("hashes", "all")
	}

	return c.postRequest(ctx, buildPath("/torrents/start"), val, nil)
}
