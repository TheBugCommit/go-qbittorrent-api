package qbittorrent

import (
	"net/http"
	"sync"
)

type Client struct {
	baseUrl    string
	httpClient *http.Client
	username   string
	password   string
	apiKey     string // By default
	mu         sync.Mutex
}
