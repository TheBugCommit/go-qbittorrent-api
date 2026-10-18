package qbittorrent

import (
	"context"

	"github.com/TheBugCommit/go-qbittorrent-api/models"
)

// SyncMainData returns the main data snapshot or partial update since the given RID.
// A zero-value SyncMainDataRequest returns a full update.
func (c *Client) SyncMainData(ctx context.Context, req models.SyncMainDataRequest) (*models.MainData, error) {
	data := &models.MainData{}
	if err := c.getRequest(ctx, buildPath("/sync/maindata"), req.ToValues(), data); err != nil {
		return nil, err
	}

	for hash, torrent := range data.Torrents {
		torrent.Hash = hash
		data.Torrents[hash] = torrent
	}

	return data, nil
}

// SyncTorrentPeers returns the peers data snapshot or partial update for a
// specific torrent since the given RID.
func (c *Client) SyncTorrentPeers(ctx context.Context, req models.SyncTorrentPeersRequest) (*models.TorrentPeersData, error) {
	peersData := &models.TorrentPeersData{}
	if err := c.getRequest(ctx, buildPath("/sync/torrentPeers"), req.ToValues(), peersData); err != nil {
		return nil, err
	}

	return peersData, nil
}
