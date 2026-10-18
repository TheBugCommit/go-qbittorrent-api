package models

import "net/url"

// SyncMainDataRequest contains the parameters for the sync/maindata endpoint.
type SyncMainDataRequest struct {
	// RID is the response ID of the last server reply. If nil or 0, a full
	// update will be returned.
	RID *int
}

// ToValues encodes the request into URL-form parameters.
func (r *SyncMainDataRequest) ToValues() url.Values {
	values := url.Values{}
	setIntParam(values, "rid", r.RID)

	return values
}

// FromValues decodes URL-form parameters into the request. It is the inverse
// of ToValues.
func (r *SyncMainDataRequest) FromValues(values url.Values) error {
	var err error
	if r.RID, err = intParam(values, "rid"); err != nil {
		return err
	}

	return nil
}

// MainData represents the response returned by /api/v2/sync/maindata.
type MainData struct {
	RID               int                 `json:"rid"`
	FullUpdate        bool                `json:"full_update"`
	Torrents          map[string]Torrent  `json:"torrents"`
	TorrentsRemoved   []string            `json:"torrents_removed"`
	Categories        map[string]Category `json:"categories"`
	CategoriesRemoved []string            `json:"categories_removed"`
	Tags              []string            `json:"tags"`
	TagsRemoved       []string            `json:"tags_removed"`
	Trackers          map[string][]string `json:"trackers"`
	TrackersRemoved   []string            `json:"trackers_removed"`
	ServerState       ServerState         `json:"server_state"`
}

// ServerState represents the global transfer info and server statistics.
type ServerState struct {
	AllTimeDL            int64  `json:"alltime_dl"`
	AllTimeUL            int64  `json:"alltime_ul"`
	AverageTimeQueue     int64  `json:"average_time_queue"`
	ConnectionStatus     string `json:"connection_status"`
	DHTNodes             int    `json:"dht_nodes"`
	DLInfoData           int64  `json:"dl_info_data"`
	DLInfoSpeed          int64  `json:"dl_info_speed"`
	DLRateLimit          int64  `json:"dl_rate_limit"`
	FreeSpaceOnDisk      int64  `json:"free_space_on_disk"`
	GlobalRatio          string `json:"global_ratio"`
	LastExternalAddress4 string `json:"last_external_address_v4"`
	LastExternalAddress6 string `json:"last_external_address_v6"`
	QueuedIOJobs         int    `json:"queued_io_jobs"`
	Queueing             bool   `json:"queueing"`
	ReadCacheHits        string `json:"read_cache_hits"`
	ReadCacheOverload    string `json:"read_cache_overload"`
	RefreshInterval      int64  `json:"refresh_interval"`
	TotalBuffersSize     int64  `json:"total_buffers_size"`
	TotalPeerConnections int    `json:"total_peer_connections"`
	TotalQueuedSize      int64  `json:"total_queued_size"`
	TotalWastedSession   int64  `json:"total_wasted_session"`
	UPInfoData           int64  `json:"up_info_data"`
	UPInfoSpeed          int64  `json:"up_info_speed"`
	UPRateLimit          int64  `json:"up_rate_limit"`
	UseAltSpeedLimits    bool   `json:"use_alt_speed_limits"`
	WriteCacheOverload   string `json:"write_cache_overload"`
}

// SyncTorrentPeersRequest contains the parameters for the sync/torrentPeers endpoint.
type SyncTorrentPeersRequest struct {
	// Hash is the torrent hash (required).
	Hash string
	// RID is the response ID of the last server reply. If nil or 0, a full
	// update will be returned.
	RID *int
}

// ToValues encodes the request into URL-form parameters.
func (r *SyncTorrentPeersRequest) ToValues() url.Values {
	values := url.Values{}
	setStringParam(values, "hash", &r.Hash)
	setIntParam(values, "rid", r.RID)

	return values
}

// FromValues decodes URL-form parameters into the request. It is the inverse
// of ToValues.
func (r *SyncTorrentPeersRequest) FromValues(values url.Values) error {
	if param := stringParam(values, "hash"); param != nil {
		r.Hash = *param
	}

	var err error
	if r.RID, err = intParam(values, "rid"); err != nil {
		return err
	}

	return nil
}

// TorrentPeersData represents the response returned by /api/v2/sync/torrentPeers.
type TorrentPeersData struct {
	FullUpdate   bool                   `json:"full_update"`
	Peers        map[string]TorrentPeer `json:"peers"`
	PeersRemoved []string               `json:"peers_removed"`
	RID          int                    `json:"rid"`
	ShowFlags    bool                   `json:"show_flags"`
}

// TorrentPeer represents a single peer connected to a torrent.
type TorrentPeer struct {
	Client        string  `json:"client"`
	Connection    string  `json:"connection"`
	Country       string  `json:"country"`
	CountryCode   string  `json:"country_code"`
	DownloadSpeed int64   `json:"dl_speed"`
	Downloaded    int64   `json:"downloaded"`
	Files         string  `json:"files"`
	Flags         string  `json:"flags"`
	FlagsDesc     string  `json:"flags_desc"`
	HostName      string  `json:"host_name"`
	IP            string  `json:"ip"`
	PeerIDClient  string  `json:"peer_id_client"`
	Port          int     `json:"port"`
	Progress      float64 `json:"progress"`
	Relevance     float64 `json:"relevance"`
	UploadSpeed   int64   `json:"up_speed"`
	Uploaded      int64   `json:"uploaded"`
}
