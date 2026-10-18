package qbittorrent

import (
	"context"
	"net/url"
	"strings"

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
	c.replaceHashesParam(&val)

	return c.postRequest(ctx, buildPath("/torrents/delete"), val, nil)
}

// PauseTorrents pause one or more torrents. Pause all torrents, without passing torrent hashes.
func (c *Client) PauseTorrents(ctx context.Context, req models.TorrentHashesRequest) error {
	val := req.ToValues()
	c.replaceHashesParam(&val)

	return c.postRequest(ctx, buildPath("/torrents/stop"), val, nil)
}

// ResumeTorrents resume one or more torrents. Resume all torrents, without passing torrent hashes.
func (c *Client) ResumeTorrents(ctx context.Context, req models.TorrentHashesRequest) error {
	val := req.ToValues()
	c.replaceHashesParam(&val)

	return c.postRequest(ctx, buildPath("/torrents/start"), val, nil)
}

// RecheckTorrents rechecks one or more torrents. Rechecks all torrents, without passing torrent hashes.
func (c *Client) RecheckTorrents(ctx context.Context, req models.TorrentHashesRequest) error {
	val := req.ToValues()
	c.replaceHashesParam(&val)

	return c.postRequest(ctx, buildPath("/torrents/recheck"), val, nil)
}

// ReannounceTorrents reannounce one or more torrents. Reannounce all torrents, without passing torrent hashes.
func (c *Client) ReannounceTorrents(ctx context.Context, req models.TorrentHashesRequest) error {
	val := req.ToValues()
	c.replaceHashesParam(&val)

	return c.postRequest(ctx, buildPath("/torrents/reannounce"), val, nil)
}

// SetTorrentsCategory sets one or more torrents category. Set all torrents category, without passing torrent hashes.
func (c *Client) SetTorrentsCategory(ctx context.Context, req models.SetTorrentsCategoryRequest) error {
	val := req.ToValues()
	c.replaceHashesParam(&val)

	return c.postRequest(ctx, buildPath("/torrents/setCategory"), val, nil)
}

// SetAutoTMM sets one or more automatic torrents management . Set all automatic torrents management, without passing torrent hashes.
func (c *Client) SetAutoTMM(ctx context.Context, req models.SetAutoTMMRequest) error {
	val := req.ToValues()
	c.replaceHashesParam(&val)

	return c.postRequest(ctx, buildPath("/torrents/setAutoManagement"), val, nil)
}

// AddTagsToTorrents adds one or more tags to torrents. Set tags to all torrents, without passing torrent hashes.
func (c *Client) AddTagsToTorrents(ctx context.Context, req models.AddTagsToTorrentsRequest) error {
	val := req.ToValues()
	c.replaceHashesParam(&val)

	return c.postRequest(ctx, buildPath("/torrents/addTags"), val, nil)
}

// RemoveTorrentTags remove one or more tags from torrents. Remove tags from all torrents, without passing torrent hashes.
func (c *Client) RemoveTorrentTags(ctx context.Context, req models.RemoveTorrentTagsRequest) error {
	val := req.ToValues()
	c.replaceHashesParam(&val)

	return c.postRequest(ctx, buildPath("/torrents/removeTags"), val, nil)
}

func (c *Client) replaceHashesParam(val *url.Values) {
	if strings.TrimSpace(val.Get("hashes")) == "" {
		val.Set("hashes", "all")
	}
}
