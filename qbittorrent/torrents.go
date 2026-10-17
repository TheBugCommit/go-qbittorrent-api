package qbittorrent

import (
	"context"
	"strings"

	"github.com/TheBugCommit/go-qbittorrent-api/models"
)

//*
//TODO: Sync,SetTorrentCategory, SetAutomaticTorrentManagment, AddTorrentTags,RemoveTorrentTags

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

// AddTorrents adds one or more torrents from URLs o File paths or magnet links.
func (c *Client) AddTorrents(ctx context.Context, req models.AddTorrentsRequest) error {
	uri := buildPath("/torrents/add")

	if len(req.Files) > 0 {
		return c.postFileRequest(ctx, uri, "torrents", req.Files, req.ToValues())

	}

	return c.postRequest(ctx, uri, req.ToValues(), nil)
}

// DeleteTorrents deletes one or more torrents. Deletes all torrents, without passing torrent hashes.
func (c *Client) DeleteTorrents(ctx context.Context, req models.DeleteTorrentsRequest) error {
	val := req.ToValues()
	if strings.TrimSpace(val.Get("hashes")) == "" {
		val.Set("hashes", "all")
	}

	return c.postRequest(ctx, buildPath("/torrents/delete"), val, nil)
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

// RecheckTorrents rechecks one or more torrents. Rechecks all torrents, without passing torrent hashes.
func (c *Client) RecheckTorrents(ctx context.Context, req models.TorrentHashesRequest) error {
	val := req.ToValues()
	if strings.TrimSpace(val.Get("hashes")) == "" {
		val.Set("hashes", "all")
	}

	return c.postRequest(ctx, buildPath("/torrents/recheck"), val, nil)
}

// ReannounceTorrents reannounce one or more torrents. Reannounce all torrents, without passing torrent hashes.
func (c *Client) ReannounceTorrents(ctx context.Context, req models.TorrentHashesRequest) error {
	val := req.ToValues()
	if strings.TrimSpace(val.Get("hashes")) == "" {
		val.Set("hashes", "all")
	}

	return c.postRequest(ctx, buildPath("/torrents/reannounce"), val, nil)
}

// SetTorrentsCategory sets one or more torrents category. Set all torrents category, without passing torrent hashes.
func (c *Client) SetTorrentsCategory(ctx context.Context, req models.SetTorrentsCategoryRequest) error {
	val := req.ToValues()
	if strings.TrimSpace(val.Get("hashes")) == "" {
		val.Set("hashes", "all")
	}

	return c.postRequest(ctx, buildPath("/torrents/setCategory"), val, nil)
}
