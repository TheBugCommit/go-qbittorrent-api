# go-qbittorrent-api
> [!WARNING]
> This library is currently under active development and the API is **not yet complete**. Breaking changes may occur in future releases until a stable $v1.0.0$ is reached. Contributions and feature requests are very welcome!

A clean, idiomatic, and robust Go client library for the **qBittorrent Web API v2**. Designed to be concurrent-safe, lightweight, and fully compatible with **qBittorrent 5.0+**.

## Features

* **qBittorrent 5.0+ Ready:** Tailored for modern qBittorrent instances utilizing JSON-based endpoints and modern authentication flows.
* **Full Authentication Support:** Log in via credentials (with smart cookie jar management and ban protection) or use API keys (`qbt_...`) seamlessly.
* **Concurrent-Safe:** Built with concurrency in mind for multi-goroutine applications.

## Implemented Endpoints

The following API groups and endpoints are currently implemented and fully functional:

| Domain | Implemented Endpoints / Features |
| :--- | :--- |
| **App** | Version, WebAPIVersion |
| **Auth** | Login, Logout |
| **Transfer** | Global transfer info, Get/Set upload & download speed limits |
| **Torrents** | Get list, properties, add (URLs, magnets, multipart files), delete, pause/resume, recheck, reannounce, categories, automatic TMM, and tag management |
| **Categories** | Get, create, edit save paths, and remove categories |
| **Tags** | Get, create, and delete tags |
| **Sync** | Sync main application data snapshot/updates and torrent peers |

## Installation

```bash
go get github.com/TheBugCommit/go-qbittorrent-api
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"

	qbittorrent "github.com/TheBugCommit/go-qbittorrent-api"
	"github.com/TheBugCommit/go-qbittorrent-api/models"
)

func main() {
	ctx := context.Background()

	// Initialize the client
	client, err := qbittorrent.New(
		"http://localhost:8080",
		qbittorrent.WithLogin("admin", "admin@123"),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Authenticate
	if err := client.Login(ctx); err != nil {
		log.Fatalf("Login failed: %v", err)
	}

	// Get qBittorrent version
	version, err := client.Version(ctx)
	if err != nil {
		log.Fatalf("Failed to get version: %v", err)
	}
	fmt.Printf("Connected to qBittorrent %s\n", version)

	// List torrents
	torrents, err := client.GetTorrents(ctx, models.TorrentListRequest{})
	if err != nil {
		log.Fatalf("Failed to fetch torrents: %v", err)
	}

	for _, t := range torrents {
		fmt.Printf("- [%s] %s (Progress: %.2f%%)\n", t.State, t.Name, t.Progress*100)
	}
}
```

## License

This project is licensed under the terms of the **GNU General Public License v3.0 (GNU GPLv3)**. See the [LICENSE](LICENSE) file for details.