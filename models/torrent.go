package models

import "net/url"

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

// TorrentHashRequest identifies a single torrent by its hash. It is used by
// endpoints that operate on one torrent, such as GetTorrentProperties.
type TorrentHashRequest struct {
	// Hash is the torrent hash (infohash v1 or v2).
	Hash string
}

// ToValues encodes the request into URL-form parameters.
func (r *TorrentHashRequest) ToValues() url.Values {
	values := url.Values{}
	values.Set("hash", r.Hash)

	return values
}

// FromValues decodes URL-form parameters into the request. It is the inverse
// of ToValues.
func (r *TorrentHashRequest) FromValues(values url.Values) error {
	r.Hash = values.Get("hash")

	return nil
}

// Tag is a tag name as returned by the Web API.
type Tag struct {
	Name string `json:"name"`
}

// CreateTagRequest creates one or more tags.
type CreateTagRequest struct {
	// Tags are the tag names to create.
	Tags []string
}

// ToValues encodes the request into URL-form parameters.
func (r *CreateTagRequest) ToValues() url.Values {
	values := url.Values{}
	setListParam(values, "tags", r.Tags, ",")

	return values
}

// FromValues decodes URL-form parameters into the request. It is the inverse
// of ToValues.
func (r *CreateTagRequest) FromValues(values url.Values) error {
	r.Tags = listParam(values, "tags", ",")

	return nil
}

// DeleteTagRequest deletes one or more tags.
type DeleteTagRequest struct {
	// Tags are the tag names to delete.
	Tags []string
}

// ToValues encodes the request into URL-form parameters.
func (r *DeleteTagRequest) ToValues() url.Values {
	values := url.Values{}
	setListParam(values, "tags", r.Tags, ",")

	return values
}

// FromValues decodes URL-form parameters into the request. It is the inverse
// of ToValues.
func (r *DeleteTagRequest) FromValues(values url.Values) error {
	r.Tags = listParam(values, "tags", ",")

	return nil
}

// ToValues encodes the request into URL-form parameters.
func (r *AddTorrentRequest) ToValues() url.Values {
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

	return values
}

// FromValues decodes URL-form parameters into the request. It is the inverse
// of ToValues.
func (r *AddTorrentRequest) FromValues(values url.Values) error {
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
