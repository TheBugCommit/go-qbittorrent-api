package models

// TransferInfo is the global transfer information returned by the Web API.
type TransferInfo struct {
	DownloadSpeed     int64  `json:"dl_info_speed"`
	DownloadData      int64  `json:"dl_info_data"`
	UploadSpeed       int64  `json:"up_info_speed"`
	UploadData        int64  `json:"up_info_data"`
	DownloadRateLimit int64  `json:"dl_rate_limit"`
	UploadRateLimit   int64  `json:"up_rate_limit"`
	DHTNodes          int64  `json:"dht_nodes"`
	ConnectionStatus  string `json:"connection_status"`
}

// SetTransferLimitRequest sets a transfer speed limit, in bytes per second.
// A limit of 0 means unlimited.
type SetTransferLimitRequest struct {
	Limit int64
}
