package models

import (
	"net/url"
)

// AddTorrentsRequest adds one or more torrents from magnet links or URLs.
//
// Optional fields are pointers so that "not set" is distinguishable from the
// zero value (e.g. Paused = false is meaningful).
type AddTorrentsRequest struct {
	// URLs is the list of torrent URLs or magnet links. At least one is
	// required.
	URLs []string
	// Files is the torrents file paths. If len(Files) > 0 then URLs is omitted.
	Files []string
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
	// AutoTMM Whether Automatic Torrent Management should be used
	AutoTMM *bool
	// SeedingTimeLimit optional since 2.8.1 Set torrent seeding time limit. Unit in minutes
	SeedingTimeLimit *int64
	// RatioLimit optional since 2.8.1 Set torrent share ratio limit
	RatioLimit *float64
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

type Torrent struct {
	Hash                     string  `json:"hash"`
	AddedOn                  int64   `json:"added_on"`
	AmountLeft               int64   `json:"amount_left"`
	AutoTMM                  bool    `json:"auto_tmm"`
	Availability             float64 `json:"availability"`
	Category                 string  `json:"category"`
	Comment                  string  `json:"comment"`
	Completed                int64   `json:"completed"`
	CompletionOn             int64   `json:"completion_on"`
	ConnectionsCount         int     `json:"connections_count"`
	ConnectionsLimit         int     `json:"connections_limit"`
	ContentPath              string  `json:"content_path"`
	CreatedBy                string  `json:"created_by"`
	CreationDate             int64   `json:"creation_date"`
	DownloadLimit            int64   `json:"dl_limit"`
	DownloadSpeed            int64   `json:"dlspeed"`
	DownloadPath             string  `json:"download_path"`
	Downloaded               int64   `json:"downloaded"`
	DownloadedSession        int64   `json:"downloaded_session"`
	ETA                      int64   `json:"eta"`
	FirstLastPiece           bool    `json:"f_l_piece_prio"`
	ForceStart               bool    `json:"force_start"`
	HasMetadata              bool    `json:"has_metadata"`
	HasOtherAnnounceError    bool    `json:"has_other_announce_error"`
	HasTrackerError          bool    `json:"has_tracker_error"`
	HasTrackerWarning        bool    `json:"has_tracker_warning"`
	InactiveSeedingTimeLimit int64   `json:"inactive_seeding_time_limit"`
	InfohashV1               string  `json:"infohash_v1"`
	InfohashV2               string  `json:"infohash_v2"`
	LastActivity             int64   `json:"last_activity"`
	MagnetURI                string  `json:"magnet_uri"`
	MaxInactiveSeedingTime   int64   `json:"max_inactive_seeding_time"`
	MaxRatio                 float64 `json:"max_ratio"`
	MaxSeedingTime           int64   `json:"max_seeding_time"`
	Name                     string  `json:"name"`
	NumComplete              int     `json:"num_complete"`
	NumIncomplete            int     `json:"num_incomplete"`
	NumLeechs                int     `json:"num_leechs"`
	NumSeeds                 int     `json:"num_seeds"`
	PieceSize                int64   `json:"piece_size"`
	PiecesHave               int     `json:"pieces_have"`
	PiecesNum                int     `json:"pieces_num"`
	Popularity               float64 `json:"popularity"`
	Priority                 int     `json:"priority"`
	Private                  bool    `json:"private"`
	Progress                 float64 `json:"progress"`
	Ratio                    float64 `json:"ratio"`
	RatioLimit               float64 `json:"ratio_limit"`
	Reannounce               int64   `json:"reannounce"`
	RootPath                 string  `json:"root_path"`
	SavePath                 string  `json:"save_path"`
	SeedingTime              int64   `json:"seeding_time"`
	SeedingTimeLimit         int64   `json:"seeding_time_limit"`
	SeenComplete             int64   `json:"seen_complete"`
	SequentialDL             bool    `json:"seq_dl"`
	ShareLimitAction         string  `json:"share_limit_action"`
	Size                     int64   `json:"size"`
	State                    string  `json:"state"`
	SuperSeeding             bool    `json:"super_seeding"`
	Tags                     string  `json:"tags"`
	TimeActive               int64   `json:"time_active"`
	TotalSize                int64   `json:"total_size"`
	TotalWasted              int64   `json:"total_wasted"`
	Tracker                  string  `json:"tracker"`
	TrackersCount            int     `json:"trackers_count"`
	UploadLimit              int64   `json:"up_limit"`
	Uploaded                 int64   `json:"uploaded"`
	UploadedSession          int64   `json:"uploaded_session"`
	UploadSpeed              int64   `json:"upspeed"`
}

// TorrentProperties is the detailed information returned for a single torrent.
type TorrentProperties struct {
	AdditionDate           int64   `json:"addition_date"`
	Availability           float64 `json:"availability"`
	Comment                string  `json:"comment"`
	CompletionDate         int64   `json:"completion_date"`
	Connections            int     `json:"nb_connections"`
	ConnectionsLimit       int     `json:"nb_connections_limit"`
	CreatedBy              string  `json:"created_by"`
	CreationDate           int64   `json:"creation_date"`
	DownloadLimit          int64   `json:"dl_limit"`
	DownloadPath           string  `json:"download_path"`
	DownloadSpeed          int64   `json:"dl_speed"`
	DownloadSpeedAvg       int64   `json:"dl_speed_avg"`
	ETA                    int64   `json:"eta"`
	HasMetadata            bool    `json:"has_metadata"`
	Hash                   string  `json:"hash"`
	InfohashV1             string  `json:"infohash_v1"`
	InfohashV2             string  `json:"infohash_v2"`
	IsPrivate              bool    `json:"is_private"`
	LastSeen               int64   `json:"last_seen"`
	Name                   string  `json:"name"`
	Peers                  int     `json:"peers"`
	PeersTotal             int     `json:"peers_total"`
	PieceSize              int64   `json:"piece_size"`
	PiecesHave             int     `json:"pieces_have"`
	PiecesNum              int     `json:"pieces_num"`
	Popularity             float64 `json:"popularity"`
	Private                bool    `json:"private"`
	Progress               float64 `json:"progress"`
	Reannounce             int64   `json:"reannounce"`
	Recheck                int64   `json:"recheck"`
	SavePath               string  `json:"save_path"`
	SeedingTime            int64   `json:"seeding_time"`
	Seeds                  int     `json:"seeds"`
	SeedsTotal             int     `json:"seeds_total"`
	ShareRatio             float64 `json:"share_ratio"`
	TimeElapsed            int64   `json:"time_elapsed"`
	TotalDownloaded        int64   `json:"total_downloaded"`
	TotalDownloadedSession int64   `json:"total_downloaded_session"`
	TotalSize              int64   `json:"total_size"`
	TotalUploaded          int64   `json:"total_uploaded"`
	TotalUploadedSession   int64   `json:"total_uploaded_session"`
	TotalWasted            int64   `json:"total_wasted"`
	UploadLimit            int64   `json:"up_limit"`
	UploadSpeed            int64   `json:"up_speed"`
	UploadSpeedAvg         int64   `json:"up_speed_avg"`
}

// TorrentHashRequest identifies a single torrent by its hash. It is used by
// endpoints that operate on one torrent, such as GetTorrentProperties.
type TorrentHashRequest struct {
	// Hash is the torrent hash (infohash v1 or v2).
	Hash string
}

type TorrentHashesRequest struct {
	Hashes []string
}

type DeleteTorrentsRequest struct {
	Hashes      []string
	DeleteFiles bool
}

type SetTorrentsCategoryRequest struct {
	Hashes   []string
	Category string
}

type SetAutoTMMRequest struct {
	Hashes []string
	Enable bool
}

type AddTagsToTorrentsRequest struct {
	Hashes []string
	Tags   []string
}

type RemoveTorrentTagsRequest struct {
	Hashes []string
	Tags   []string
}

// ToValues encodes the request into URL-form parameters.
func (r *RemoveTorrentTagsRequest) ToValues() url.Values {
	values := url.Values{}
	setListParam(values, "hashes", r.Hashes, "|")
	setListParam(values, "tags", r.Tags, ",")

	return values
}

// FromValues decodes URL-form parameters into the request.
func (r *RemoveTorrentTagsRequest) FromValues(values url.Values) error {
	r.Hashes = listParam(values, "hashes", "|")
	r.Tags = listParam(values, "tags", ",")

	return nil
}

// ToValues encodes the request into URL-form parameters.
func (r *AddTagsToTorrentsRequest) ToValues() url.Values {
	values := url.Values{}
	setListParam(values, "hashes", r.Hashes, "|")
	setListParam(values, "tags", r.Tags, ",")

	return values
}

// FromValues decodes URL-form parameters into the request.
func (r *AddTagsToTorrentsRequest) FromValues(values url.Values) error {
	r.Hashes = listParam(values, "hashes", "|")
	r.Tags = listParam(values, "tags", ",")

	return nil
}

// ToValues encodes the request into URL-form parameters.
func (r *SetAutoTMMRequest) ToValues() url.Values {
	values := url.Values{}
	setListParam(values, "hashes", r.Hashes, "|")
	setBoolParam(values, "enable", &r.Enable)

	return values
}

// FromValues decodes URL-form parameters into the request.
func (r *SetAutoTMMRequest) FromValues(values url.Values) error {
	r.Hashes = listParam(values, "hashes", "|")

	val, err := boolParam(values, "enable")
	if err != nil {
		return err
	}
	if val != nil {
		r.Enable = *val
	} else {
		r.Enable = false
	}

	return nil
}

// ToValues encodes the request into URL-form parameters.
func (r *SetTorrentsCategoryRequest) ToValues() url.Values {
	values := url.Values{}
	setListParam(values, "hashes", r.Hashes, "|")
	setStringParam(values, "category", &r.Category)

	return values
}

// FromValues decodes URL-form parameters into the request.
func (r *SetTorrentsCategoryRequest) FromValues(values url.Values) error {
	r.Hashes = listParam(values, "hashes", "|")
	if param := stringParam(values, "category"); param != nil {
		r.Category = *param
	}

	return nil
}

// ToValues encodes the request into URL-form parameters.
func (r *DeleteTorrentsRequest) ToValues() url.Values {
	values := url.Values{}
	setListParam(values, "hashes", r.Hashes, "|")
	setBoolParam(values, "deleteFiles", &r.DeleteFiles)

	return values
}

// FromValues decodes URL-form parameters into the request.
func (r *DeleteTorrentsRequest) FromValues(values url.Values) error {
	r.Hashes = listParam(values, "hashes", "|")

	val, err := boolParam(values, "deleteFiles")
	if err != nil {
		return err
	}
	if val != nil {
		r.DeleteFiles = *val
	} else {
		r.DeleteFiles = false
	}

	return nil
}

// ToValues encodes the request into URL-form parameters.
func (r *TorrentHashesRequest) ToValues() url.Values {
	values := url.Values{}
	setListParam(values, "hashes", r.Hashes, "|")

	return values
}

// FromValues decodes URL-form parameters into the request. It is the inverse
// of ToValues.
func (r *TorrentHashesRequest) FromValues(values url.Values) error {
	r.Hashes = listParam(values, "hashes", "|")

	return nil
}

// ToValues encodes the request into URL-form parameters.
func (r *TorrentHashRequest) ToValues() url.Values {
	values := url.Values{}
	setStringParam(values, "hash", &r.Hash)

	return values
}

// FromValues decodes URL-form parameters into the request. It is the inverse
// of ToValues.
func (r *TorrentHashRequest) FromValues(values url.Values) error {
	if param := stringParam(values, "hash"); param != nil {
		r.Hash = *param
	}

	return nil
}

// ToValues encodes the request into URL-form parameters.
func (r *AddTorrentsRequest) ToValues() url.Values {
	values := url.Values{}
	setListParam(values, "urls", r.URLs, "\n")
	setStringParam(values, "savepath", r.SavePath)
	setStringParam(values, "category", r.Category)
	setListParam(values, "tags", r.Tags, ",")
	setBoolParam(values, "skip_checking", r.SkipChecking)
	setBoolParam(values, "paused", r.Paused)
	setBoolParam(values, "root_folder", r.RootFolder)
	setStringParam(values, "rename", r.Rename)
	setInt64Param(values, "dlLimit", r.DownloadLimit)
	setInt64Param(values, "upLimit", r.UploadLimit)
	setBoolParam(values, "sequentialDownload", r.SequentialDownload)
	setBoolParam(values, "firstLastPiecePrio", r.FirstLastPiecePrio)
	setBoolParam(values, "autoTMM", r.AutoTMM)
	setInt64Param(values, "seedingTimeLimit", r.SeedingTimeLimit)
	setFloat64Param(values, "ratioLimit", r.RatioLimit)

	return values
}

// FromValues decodes URL-form parameters into the request. It is the inverse
// of ToValues.
func (r *AddTorrentsRequest) FromValues(values url.Values) error {
	r.URLs = listParam(values, "urls", "\n")
	r.SavePath = stringParam(values, "savepath")
	r.Category = stringParam(values, "category")
	r.Tags = listParam(values, "tags", ",")
	r.Rename = stringParam(values, "rename")

	var err error
	if r.SkipChecking, err = boolParam(values, "skip_checking"); err != nil {
		return err
	}
	if r.Paused, err = boolParam(values, "paused"); err != nil {
		return err
	}
	if r.RootFolder, err = boolParam(values, "root_folder"); err != nil {
		return err
	}
	if r.DownloadLimit, err = int64Param(values, "dlLimit"); err != nil {
		return err
	}
	if r.UploadLimit, err = int64Param(values, "upLimit"); err != nil {
		return err
	}
	if r.SequentialDownload, err = boolParam(values, "sequentialDownload"); err != nil {
		return err
	}
	if r.FirstLastPiecePrio, err = boolParam(values, "firstLastPiecePrio"); err != nil {
		return err
	}
	if r.AutoTMM, err = boolParam(values, "autoTMM"); err != nil {
		return err
	}
	if r.SeedingTimeLimit, err = int64Param(values, "seedingTimeLimit"); err != nil {
		return err
	}
	if r.RatioLimit, err = float64Param(values, "ratioLimit"); err != nil {
		return err
	}

	return nil
}

// ToValues encodes the request into URL-form parameters.
func (r *TorrentListRequest) ToValues() url.Values {
	values := url.Values{}
	setStringParam(values, "filter", r.Filter)
	setStringParam(values, "category", r.Category)
	setStringParam(values, "tag", r.Tag)
	setListParam(values, "hashes", r.Hashes, "|")
	setStringParam(values, "sort", r.Sort)
	setBoolParam(values, "reverse", r.Reverse)
	setIntParam(values, "limit", r.Limit)
	setIntParam(values, "offset", r.Offset)

	return values
}

// FromValues decodes URL-form parameters into the request. It is the inverse
// of ToValues.
func (r *TorrentListRequest) FromValues(values url.Values) error {
	r.Filter = stringParam(values, "filter")
	r.Category = stringParam(values, "category")
	r.Tag = stringParam(values, "tag")
	r.Hashes = listParam(values, "hashes", "|")
	r.Sort = stringParam(values, "sort")

	var err error
	if r.Reverse, err = boolParam(values, "reverse"); err != nil {
		return err
	}
	if r.Limit, err = intParam(values, "limit"); err != nil {
		return err
	}
	if r.Offset, err = intParam(values, "offset"); err != nil {
		return err
	}

	return nil
}
