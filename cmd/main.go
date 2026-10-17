// Command main exercises every request currently supported by the
// qbittorrent client against a live qBittorrent instance.
//
// Connection settings are read from the environment:
//
//	QB_URL       base URL, default http://localhost:8080
//	QB_API_KEY   API key; when set, username/password are ignored
//	QB_USERNAME  WebUI username, default admin
//	QB_PASSWORD  WebUI password, default adminadmin
//	QB_REFERER   optional Referer header override
//
// Read-only requests always run. Requests that change server state run only
// when the -write flag is given, and AddTorrent additionally needs -torrent.
//
//	go run ./cmd -write -torrent 'magnet:?xt=urn:btih:...'
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TheBugCommit/go-qbittorrent-api/models"
	"github.com/TheBugCommit/go-qbittorrent-api/qbittorrent"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	write := flag.Bool("write", false, "run requests that change server state")
	torrent := flag.String("torrent", "", "magnet link or URL to add when -write is set")
	category := flag.String("category", "go-qbittorrent-api-test", "category name used by the write operations")
	tag := flag.String("tag", "go-qbittorrent-api-test", "tag name used by the write operations")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	client, err := newClient()
	if err != nil {
		return err
	}
	defer func() { _ = client.Close() }()

	authMode := "username/password"
	if os.Getenv("QB_API_KEY") != "" {
		authMode = "API key"
	}
	fmt.Printf("base URL : %s\n", client.BaseURL())
	fmt.Printf("referer  : %s\n", client.Referer())
	fmt.Printf("auth     : %s\n\n", authMode)

	fmt.Println("--- authentication ---")
	if err := client.Login(ctx); err != nil {
		fmt.Printf("  ✗ %-20s %v\n", "Login", err)

		return fmt.Errorf("login failed, stopping: %w", err)
	}
	fmt.Printf("  ✓ %-20s ok\n", "Login")

	fmt.Println("\n--- app ---")
	report(ctx, "Version", client.Version)
	report(ctx, "WebAPIVersion", client.WebAPIVersion)

	fmt.Println("\n--- categories ---")
	if cats, err := client.GetCategories(ctx); err != nil {
		fmt.Printf("  ✗ %-20s %v\n", "GetCategories", err)
	} else {
		fmt.Printf("  ✓ %-20s %d categories\n", "GetCategories", len(cats))
		for name := range cats {
			fmt.Printf("      • %s\n", name)
		}
	}

	fmt.Println("\n--- tags ---")
	if tags, err := client.GetTags(ctx); err != nil {
		fmt.Printf("  ✗ %-20s %v\n", "GetTags", err)
	} else {
		fmt.Printf("  ✓ %-20s %d tags\n", "GetTags", len(tags))
		for _, tag := range tags {
			fmt.Printf("      • %s\n", tag.Name)
		}
	}

	fmt.Println("\n--- transfer ---")
	report(ctx, "GetTransferInfo", client.GetTransferInfo)
	report(ctx, "GetDownloadLimit", client.GetDownloadLimit)
	report(ctx, "GetUploadLimit", client.GetUploadLimit)

	fmt.Println("\n--- torrents ---")
	torrents, err := client.GetTorrents(ctx, models.TorrentListRequest{})
	if err != nil {
		fmt.Printf("  ✗ %-20s %v\n", "GetTorrents", err)
	}
	if err == nil {
		fmt.Printf("  ✓ %-20s %d torrents\n", "GetTorrents", len(torrents))
		for _, t := range torrents {
			fmt.Printf("      • %.8s  %-10s  %s\n", t.Hash, t.State, t.Name)
		}
	}

	if len(torrents) > 0 {
		hash := torrents[0].Hash
		report(ctx, "GetTorrentProperties", func(ctx context.Context) (*models.TorrentProperties, error) {
			return client.GetTorrentProperties(ctx, hash)
		})
	} else {
		fmt.Printf("  - %-20s skipped (no torrents)\n", "GetTorrentProperties")
	}

	reportErr(ctx, "PauseTorrents", func(ctx context.Context) error {
		return client.PauseTorrents(ctx, models.TorrentHashesRequest{Hashes: []string{}})
	})

	reportErr(ctx, "ResumeTorrents", func(ctx context.Context) error {
		return client.ResumeTorrents(ctx, models.TorrentHashesRequest{Hashes: []string{}})
	})

	if !*write {
		fmt.Println("\nwrite operations skipped, pass -write to run them")
		return nil
	}

	fmt.Println("\n--- write operations (state changing) ---")

	reportErr(ctx, "CreateCategory", func(ctx context.Context) error {
		return client.CreateCategory(ctx, models.CreateCategoryRequest{
			Name:     *category,
			SavePath: "/tmp/go-qbittorrent-api-test",
		})
	})
	reportErr(ctx, "EditCategory", func(ctx context.Context) error {
		return client.EditCategory(ctx, models.EditCategoryRequest{
			Name:     *category,
			SavePath: "/tmp/go-qbittorrent-api-test-edited",
		})
	})
	reportErr(ctx, "RemoveCategories", func(ctx context.Context) error {
		return client.RemoveCategories(ctx, models.RemoveCategoriesRequest{
			Names: []string{*category},
		})
	})

	// The limits are set back to their current value, so the endpoint is
	// exercised without changing the server configuration.
	if current, err := client.GetDownloadLimit(ctx); err == nil {
		reportErr(ctx, "SetDownloadLimit", func(ctx context.Context) error {
			return client.SetDownloadLimit(ctx, models.SetTransferLimitRequest{Limit: current})
		})
	}
	if current, err := client.GetUploadLimit(ctx); err == nil {
		reportErr(ctx, "SetUploadLimit", func(ctx context.Context) error {
			return client.SetUploadLimit(ctx, models.SetTransferLimitRequest{Limit: current})
		})
	}

	if *torrent == "" {
		fmt.Printf("  - %-20s skipped (pass -torrent to add one)\n", "AddTorrent")
	} else {
		fmt.Printf("  ! %-20s adding %q (this starts a real download)\n", "AddTorrent", *torrent)
		/*reportErr(ctx, "AddTorrent", func(ctx context.Context) error {
			path := "/data/downloads"
			return client.AddTorrent(ctx, models.AddTorrentRequest{
				URLs:     []string{*torrent},
				SavePath: &path,
			})
		})*/

		reportErr(ctx, "AddTorrents", func(ctx context.Context) error {
			path := "/data/downloads"
			autoTMM := false
			return client.AddTorrents(ctx, models.AddTorrentsRequest{
				Files:    []string{*torrent},
				SavePath: &path,
				AutoTMM:  &autoTMM,
			})
		})

		reportErr(ctx, "DeleteTorrents", func(ctx context.Context) error {
			return client.DeleteTorrents(ctx, models.DeleteTorrentsRequest{
				DeleteFiles: false,
			})
		})
	}

	// Tags
	reportErr(ctx, "CreateTags", func(ctx context.Context) error {
		return client.CreateTags(ctx, models.CreateTagRequest{Tags: []string{*tag}})
	})

	reportErr(ctx, "DeleteTags", func(ctx context.Context) error {
		return client.CreateTags(ctx, models.CreateTagRequest{Tags: []string{*tag}})
	})

	reportErr(ctx, "Logout", client.Logout)

	return nil
}

// newClient builds a client from the environment variables described above.
func newClient() (*qbittorrent.Client, error) {
	baseURL := envOr("QB_URL", "http://localhost:8080")

	opts := []qbittorrent.Option{qbittorrent.WithTimeout(15 * time.Second)}

	if key := os.Getenv("QB_API_KEY"); key != "" {
		opts = append(opts, qbittorrent.WithAPIKey(key))
	} else {
		opts = append(opts, qbittorrent.WithLogin(
			envOr("QB_USERNAME", "admin"),
			envOr("QB_PASSWORD", "e1d76ae416763235bf2509ea4d7e5289d834ac71a2161dbb"),
		))
	}

	if referer := os.Getenv("QB_REFERER"); referer != "" {
		opts = append(opts, qbittorrent.WithReferer(referer))
	}

	return qbittorrent.New(baseURL, opts...)
}

func envOr(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

// report runs a request returning a value and prints its outcome.
func report[T any](ctx context.Context, name string, fn func(context.Context) (T, error)) {
	value, err := fn(ctx)
	if err != nil {
		fmt.Printf("  ✗ %-20s %v\n", name, err)
		return
	}

	fmt.Printf("  ✓ %-20s %+v\n", name, value)
}

// reportErr runs a request returning only an error and prints its outcome.
func reportErr(ctx context.Context, name string, fn func(context.Context) error) {
	if err := fn(ctx); err != nil {
		fmt.Printf("  ✗ %-20s %v\n", name, err)
		return
	}

	fmt.Printf("  ✓ %-20s ok\n", name)
}
