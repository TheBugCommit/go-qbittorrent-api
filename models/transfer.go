package models

import (
	"net/url"
	"strconv"
)

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

// ToValues encodes the request into URL-form parameters.
func (r *SetTransferLimitRequest) ToValues() url.Values {
	values := url.Values{}
	values.Set("limit", strconv.FormatInt(r.Limit, 10))

	return values
}

// FromValues decodes URL-form parameters into the request. It is the inverse
// of ToValues.
func (r *SetTransferLimitRequest) FromValues(values url.Values) error {
	limit, err := int64Param(values, "limit")
	if err != nil {
		return err
	}
	if limit != nil {
		r.Limit = *limit
	}

	return nil
}
