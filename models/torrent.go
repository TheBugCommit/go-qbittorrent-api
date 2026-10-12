package models

// AddTorrentRequest adds one or more torrents from magnet links or URLs.
//
// Optional fields are pointers so that "not set" is distinguishable from the
// zero value (e.g. Paused = false is meaningful).
//
// Uploading a .torrent file is not supported by this request; only URLs and
// magnet links are sent.
type AddTorrentRequest struct {
	// URLs is the list of torrent URLs or magnet links. At least one is
	// required.
	URLs []string
	// SavePath is the download location for the torrent content.
	SavePath *string
	// Category is the category to assign.
	Category *string
	// Tags are the tags to assign.
	Tags []string
	// SkipChecking skips the hash check of the torrent.
	SkipChecking *bool
	// Paused adds the torrent in a paused state.
	Paused *bool
	// RootFolder creates the root folder even for single-file torrents.
	RootFolder *bool
	// Rename renames the torrent.
	Rename *string
	// DownloadLimit is the download speed limit in bytes per second.
	DownloadLimit *int64
	// UploadLimit is the upload speed limit in bytes per second.
	UploadLimit *int64
	// SequentialDownload enables sequential downloading.
	SequentialDownload *bool
	// FirstLastPiecePrio prioritises the first and last pieces.
	FirstLastPiecePrio *bool
}

// TorrentListRequest filters the torrents returned by GetTorrents. All fields
// are optional; a zero value returns every torrent.
type TorrentListRequest struct {
	// Filter restricts the result, e.g. "downloading", "seeding",
	// "completed", "stopped", "running", "stalled", "errored".
	Filter *string
	// Category restricts the result to a category name.
	Category *string
	// Tag restricts the result to a tag.
	Tag *string
	// Hashes restricts the result to the given torrent hashes.
	Hashes []string
	// Sort orders the result, e.g. "name", "size", "progress".
	Sort *string
	// Reverse reverses the sort order.
	Reverse *bool
	// Limit is the maximum number of torrents to return.
	Limit *int
	// Offset skips the first Offset torrents.
	Offset *int
}

// Torrent is a torrent as returned by the torrent list endpoint.
type Torrent struct {
	Hash           string  `json:"hash"`
	Name           string  `json:"name"`
	Size           int64   `json:"size"`
	TotalSize      int64   `json:"total_size"`
	Progress       float64 `json:"progress"`
	DownloadSpeed  int64   `json:"dlspeed"`
	UploadSpeed    int64   `json:"upspeed"`
	Downloaded     int64   `json:"downloaded"`
	Uploaded       int64   `json:"uploaded"`
	DownloadLimit  int64   `json:"dl_limit"`
	UploadLimit    int64   `json:"up_limit"`
	AmountLeft     int64   `json:"amount_left"`
	Completed      int64   `json:"completed"`
	ETA            int64   `json:"eta"`
	Ratio          float64 `json:"ratio"`
	State          string  `json:"state"`
	Category       string  `json:"category"`
	Tags           string  `json:"tags"`
	Tracker        string  `json:"tracker"`
	SavePath       string  `json:"save_path"`
	ContentPath    string  `json:"content_path"`
	MagnetURI      string  `json:"magnet_uri"`
	InfohashV1     string  `json:"infohash_v1"`
	InfohashV2     string  `json:"infohash_v2"`
	AddedOn        int64   `json:"added_on"`
	CompletionOn   int64   `json:"completion_on"`
	SeedingTime    int64   `json:"seeding_time"`
	TimeActive     int64   `json:"time_active"`
	NumSeeds       int     `json:"num_seeds"`
	NumLeechs      int     `json:"num_leechs"`
	NumComplete    int     `json:"num_complete"`
	NumIncomplete  int     `json:"num_incomplete"`
	Priority       int     `json:"priority"`
	ForceStart     bool    `json:"force_start"`
	SequentialDL   bool    `json:"seq_dl"`
	FirstLastPiece bool    `json:"f_l_piece_prio"`
	SuperSeeding   bool    `json:"super_seeding"`
}

// TorrentProperties is the detailed information returned for a single torrent.
type TorrentProperties struct {
	SavePath               string  `json:"save_path"`
	CreationDate           int64   `json:"creation_date"`
	AdditionDate           int64   `json:"addition_date"`
	CompletionDate         int64   `json:"completion_date"`
	PieceSize              int64   `json:"piece_size"`
	Comment                string  `json:"comment"`
	CreatedBy              string  `json:"created_by"`
	TotalWasted            int64   `json:"total_wasted"`
	TotalUploaded          int64   `json:"total_uploaded"`
	TotalUploadedSession   int64   `json:"total_uploaded_session"`
	TotalDownloaded        int64   `json:"total_downloaded"`
	TotalDownloadedSession int64   `json:"total_downloaded_session"`
	TotalSize              int64   `json:"total_size"`
	UploadLimit            int64   `json:"up_limit"`
	DownloadLimit          int64   `json:"dl_limit"`
	TimeElapsed            int64   `json:"time_elapsed"`
	SeedingTime            int64   `json:"seeding_time"`
	Connections            int     `json:"nb_connections"`
	ConnectionsLimit       int     `json:"nb_connections_limit"`
	ShareRatio             float64 `json:"share_ratio"`
	DownloadSpeedAvg       int64   `json:"dl_speed_avg"`
	DownloadSpeed          int64   `json:"dl_speed"`
	UploadSpeedAvg         int64   `json:"up_speed_avg"`
	UploadSpeed            int64   `json:"up_speed"`
	ETA                    int64   `json:"eta"`
	Peers                  int     `json:"peers"`
	PeersTotal             int     `json:"peers_total"`
	PiecesHave             int     `json:"pieces_have"`
	PiecesNum              int     `json:"pieces_num"`
	Recheck                int64   `json:"recheck"`
	Seeds                  int     `json:"seeds"`
	SeedsTotal             int     `json:"seeds_total"`
	LastSeen               int64   `json:"last_seen"`
	Popularity             float64 `json:"popularity"`
	DownloadPath           string  `json:"download_path"`
	InfohashV1             string  `json:"infohash_v1"`
	InfohashV2             string  `json:"infohash_v2"`
}
