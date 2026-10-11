package qbittorrent

import (
	"net/http"
	"sync"
)

// Client is a qBittorrent Web API client. Create one with NewClient. A Client
// is intended to be safe for concurrent use by multiple goroutines.
type Client struct {
	baseURL    string
	httpClient *http.Client
	username   string
	password   string
	apiKey     string // If set, requests are authenticated with "Authorization: Bearer".
	mu         sync.Mutex
}
