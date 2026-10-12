package qbittorrent

import (
	"context"
	"net/url"
	"strconv"
	"strings"

	"github.com/TheBugCommit/go-qbittorrent-api/models"
)

// GetTorrents returns the torrents matching the given filters. A zero-value
// TorrentListRequest returns every torrent.
func (c *Client) GetTorrents(ctx context.Context, req models.TorrentListRequest) ([]models.Torrent, error) {
	var torrents []models.Torrent
	if err := c.getRequest(ctx, buildPath("/torrents/info"), torrentListValues(req), &torrents); err != nil {
		return nil, err
	}

	return torrents, nil
}

// GetTorrentProperties returns the detailed properties of a single torrent.
func (c *Client) GetTorrentProperties(ctx context.Context, hash string) (*models.TorrentProperties, error) {
	params := url.Values{"hash": {hash}}

	props := &models.TorrentProperties{}
	if err := c.getRequest(ctx, buildPath("/torrents/properties"), params, props); err != nil {
		return nil, err
	}

	return props, nil
}

// AddTorrent adds one or more torrents from URLs or magnet links. Uploading a
// .torrent file is not supported yet.
func (c *Client) AddTorrent(ctx context.Context, req models.AddTorrentRequest) error {
	return c.postRequest(ctx, buildPath("/torrents/add"), addTorrentValues(req), nil)
}

func torrentListValues(req models.TorrentListRequest) url.Values {
	v := url.Values{}

	if req.Filter != nil {
		v.Set("filter", *req.Filter)
	}
	if req.Category != nil {
		v.Set("category", *req.Category)
	}
	if req.Tag != nil {
		v.Set("tag", *req.Tag)
	}
	if len(req.Hashes) > 0 {
		v.Set("hashes", strings.Join(req.Hashes, "\n"))
	}
	if req.Sort != nil {
		v.Set("sort", *req.Sort)
	}
	if req.Reverse != nil {
		v.Set("reverse", strconv.FormatBool(*req.Reverse))
	}
	if req.Limit != nil {
		v.Set("limit", strconv.Itoa(*req.Limit))
	}
	if req.Offset != nil {
		v.Set("offset", strconv.Itoa(*req.Offset))
	}

	return v
}

func addTorrentValues(req models.AddTorrentRequest) url.Values {
	v := url.Values{}

	if len(req.URLs) > 0 {
		v.Set("urls", strings.Join(req.URLs, "\n"))
	}
	if req.SavePath != nil {
		v.Set("savepath", *req.SavePath)
	}
	if req.Category != nil {
		v.Set("category", *req.Category)
	}
	if len(req.Tags) > 0 {
		v.Set("tags", strings.Join(req.Tags, ","))
	}
	if req.SkipChecking != nil {
		v.Set("skip_checking", strconv.FormatBool(*req.SkipChecking))
	}
	if req.Paused != nil {
		v.Set("paused", strconv.FormatBool(*req.Paused))
	}
	if req.RootFolder != nil {
		v.Set("root_folder", strconv.FormatBool(*req.RootFolder))
	}
	if req.Rename != nil {
		v.Set("rename", *req.Rename)
	}
	if req.DownloadLimit != nil {
		v.Set("dlLimit", strconv.FormatInt(*req.DownloadLimit, 10))
	}
	if req.UploadLimit != nil {
		v.Set("upLimit", strconv.FormatInt(*req.UploadLimit, 10))
	}
	if req.SequentialDownload != nil {
		v.Set("sequentialDownload", strconv.FormatBool(*req.SequentialDownload))
	}
	if req.FirstLastPiecePrio != nil {
		v.Set("firstLastPiecePrio", strconv.FormatBool(*req.FirstLastPiecePrio))
	}

	return v
}
