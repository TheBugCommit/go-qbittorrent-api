package qbittorrent_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/TheBugCommit/go-qbittorrent-api/models"
	"github.com/TheBugCommit/go-qbittorrent-api/qbittorrent"
)

var (
	flagWrite    = flag.Bool("write", false, "run requests that change server state")
	flagTorrent  = flag.String("torrent", "", "torrent file path, magnet link or URL to add when -write is set")
	flagCategory = flag.String("category", "go-qbittorrent-api-test", "category name used by write operations")
	flagTag      = flag.String("tag", "go-qbittorrent-api-test", "tag name used by write operations")
)

func newTestClient(t *testing.T) *qbittorrent.Client {
	t.Helper()

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

	client, err := qbittorrent.New(baseURL, opts...)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	t.Cleanup(func() {
		_ = client.Close()
	})

	return client
}

func envOr(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func TestIntegration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	client := newTestClient(t)
	t.Logf("base URL: %s | referer: %s", client.BaseURL(), client.Referer())

	// 1. Authentication
	if err := client.Login(ctx); err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	t.Cleanup(func() {
		if *flagWrite {
			if err := client.Logout(context.Background()); err != nil {
				t.Errorf("Logout failed: %v", err)
			}
		}
	})

	// 2. App Endpoints
	t.Run("App", func(t *testing.T) {
		t.Run("Version", func(t *testing.T) {
			v, err := client.Version(ctx)
			if err != nil {
				t.Fatalf("Version() error = %v", err)
			}
			t.Logf("Version: %s", v)
		})

		t.Run("WebAPIVersion", func(t *testing.T) {
			v, err := client.WebAPIVersion(ctx)
			if err != nil {
				t.Fatalf("WebAPIVersion() error = %v", err)
			}
			t.Logf("WebAPIVersion: %s", v)
		})
	})

	// 3. Sync Endpoints
	var firstHash string
	t.Run("Sync", func(t *testing.T) {
		t.Run("SyncMainData", func(t *testing.T) {
			mainData, err := client.SyncMainData(ctx, models.SyncMainDataRequest{})
			if err != nil {
				t.Fatalf("SyncMainData() error = %v", err)
			}
			if !mainData.FullUpdate {
				t.Errorf("expected FullUpdate = true on first request, got false")
			}
			t.Logf("RID: %d, Torrents: %d, Categories: %d, Tags: %d",
				mainData.RID, len(mainData.Torrents), len(mainData.Categories), len(mainData.Tags))

			for h := range mainData.Torrents {
				firstHash = h
				break
			}

			// Test incremental update passing previous RID
			rid := mainData.RID
			partialData, err := client.SyncMainData(ctx, models.SyncMainDataRequest{RID: &rid})
			if err != nil {
				t.Fatalf("SyncMainData(RID=%d) error = %v", rid, err)
			}
			if partialData.FullUpdate {
				t.Errorf("expected FullUpdate = false on incremental request, got true")
			}
		})

		t.Run("SyncTorrentPeers", func(t *testing.T) {
			if firstHash == "" {
				t.Skip("skipping SyncTorrentPeers: no torrents available on server")
			}
			peersData, err := client.SyncTorrentPeers(ctx, models.SyncTorrentPeersRequest{
				Hash: firstHash,
			})
			if err != nil {
				t.Fatalf("SyncTorrentPeers(%s) error = %v", firstHash, err)
			}
			t.Logf("Peers for %s: %d (RID: %d)", firstHash, len(peersData.Peers), peersData.RID)
		})
	})

	// 4. Read-Only Categories, Tags & Transfer
	t.Run("Categories_Read", func(t *testing.T) {
		cats, err := client.GetCategories(ctx)
		if err != nil {
			t.Fatalf("GetCategories() error = %v", err)
		}
		t.Logf("Fetched %d categories", len(cats))
	})

	t.Run("Tags_Read", func(t *testing.T) {
		tags, err := client.GetTags(ctx)
		if err != nil {
			t.Fatalf("GetTags() error = %v", err)
		}
		t.Logf("Fetched %d tags", len(tags))
	})

	t.Run("Transfer_Read", func(t *testing.T) {
		info, err := client.GetTransferInfo(ctx)
		if err != nil {
			t.Fatalf("GetTransferInfo() error = %v", err)
		}
		t.Logf("Transfer Info: %+v", info)

		dlLimit, err := client.GetDownloadLimit(ctx)
		if err != nil {
			t.Fatalf("GetDownloadLimit() error = %v", err)
		}

		ulLimit, err := client.GetUploadLimit(ctx)
		if err != nil {
			t.Fatalf("GetUploadLimit() error = %v", err)
		}
		t.Logf("Limits -> DL: %d, UL: %d", dlLimit, ulLimit)
	})

	// 5. Read-Only Torrents
	t.Run("Torrents_Read", func(t *testing.T) {
		torrents, err := client.GetTorrents(ctx, models.TorrentListRequest{})
		if err != nil {
			t.Fatalf("GetTorrents() error = %v", err)
		}
		t.Logf("Fetched %d torrents", len(torrents))

		if len(torrents) == 0 {
			t.Skip("skipping GetTorrentProperties: no torrents available")
		}

		props, err := client.GetTorrentProperties(ctx, torrents[0].Hash)
		if err != nil {
			t.Fatalf("GetTorrentProperties(%s) error = %v", torrents[0].Hash, err)
		}
		t.Logf("Properties for %s: SavePath=%s, PieceSize=%d", torrents[0].Name, props.SavePath, props.PieceSize)
	})

	// 6. Write Operations (State-Changing)
	t.Run("WriteOperations", func(t *testing.T) {
		if !*flagWrite {
			t.Skip("write operations skipped; pass -write to run them")
		}

		t.Run("Categories_Write", func(t *testing.T) {
			tempCat := fmt.Sprintf("%s-temp", *flagCategory)

			t.Cleanup(func() {
				_ = client.RemoveCategories(ctx, models.RemoveCategoriesRequest{
					Names: []string{tempCat},
				})
			})

			if err := client.CreateCategory(ctx, models.CreateCategoryRequest{
				Name:     tempCat,
				SavePath: "/tmp/go-qbittorrent-api-test",
			}); err != nil {
				t.Fatalf("CreateCategory() error = %v", err)
			}

			if err := client.EditCategory(ctx, models.EditCategoryRequest{
				Name:     tempCat,
				SavePath: "/tmp/go-qbittorrent-api-test-edited",
			}); err != nil {
				t.Errorf("EditCategory() error = %v", err)
			}

			_ = client.CreateCategory(ctx, models.CreateCategoryRequest{
				Name:     *flagCategory,
				SavePath: "/tmp/go-qbittorrent-api-test-edited",
			})
		})

		t.Run("Tags_Write", func(t *testing.T) {
			if err := client.CreateTags(ctx, models.CreateTagRequest{
				Tags: []string{*flagTag},
			}); err != nil {
				t.Errorf("CreateTags() error = %v", err)
			}
		})

		t.Run("TransferLimits_Write", func(t *testing.T) {
			dlCurrent, err := client.GetDownloadLimit(ctx)
			if err != nil {
				t.Fatalf("GetDownloadLimit() error = %v", err)
			}
			if err := client.SetDownloadLimit(ctx, models.SetTransferLimitRequest{Limit: dlCurrent}); err != nil {
				t.Errorf("SetDownloadLimit() error = %v", err)
			}

			ulCurrent, err := client.GetUploadLimit(ctx)
			if err != nil {
				t.Fatalf("GetUploadLimit() error = %v", err)
			}
			if err := client.SetUploadLimit(ctx, models.SetTransferLimitRequest{Limit: ulCurrent}); err != nil {
				t.Errorf("SetUploadLimit() error = %v", err)
			}
		})

		t.Run("TorrentActions_Write", func(t *testing.T) {
			// If firstHash is available, operate only on that torrent instead of "all"
			var targetHashes []string
			if firstHash != "" {
				targetHashes = []string{firstHash}
			}

			if err := client.PauseTorrents(ctx, models.TorrentHashesRequest{Hashes: targetHashes}); err != nil {
				t.Errorf("PauseTorrents() error = %v", err)
			}
			if err := client.ResumeTorrents(ctx, models.TorrentHashesRequest{Hashes: targetHashes}); err != nil {
				t.Errorf("ResumeTorrents() error = %v", err)
			}
			if err := client.RecheckTorrents(ctx, models.TorrentHashesRequest{Hashes: targetHashes}); err != nil {
				t.Errorf("RecheckTorrents() error = %v", err)
			}
			if err := client.ReannounceTorrents(ctx, models.TorrentHashesRequest{Hashes: targetHashes}); err != nil {
				t.Errorf("ReannounceTorrents() error = %v", err)
			}
			if err := client.SetTorrentsCategory(ctx, models.SetTorrentsCategoryRequest{
				Hashes:   targetHashes,
				Category: *flagCategory,
			}); err != nil {
				t.Errorf("SetTorrentsCategory() error = %v", err)
			}
			if err := client.SetAutoTMM(ctx, models.SetAutoTMMRequest{
				Hashes: targetHashes,
				Enable: true,
			}); err != nil {
				t.Errorf("SetAutoTMM() error = %v", err)
			}
			if err := client.AddTagsToTorrents(ctx, models.AddTagsToTorrentsRequest{
				Hashes: targetHashes,
				Tags:   []string{*flagTag},
			}); err != nil {
				t.Errorf("AddTagsToTorrents() error = %v", err)
			}
			if err := client.RemoveTorrentTags(ctx, models.RemoveTorrentTagsRequest{
				Hashes: targetHashes,
				Tags:   []string{*flagTag},
			}); err != nil {
				t.Errorf("RemoveTorrentTags() error = %v", err)
			}
		})

		t.Run("AddTorrents_Write", func(t *testing.T) {
			if *flagTorrent == "" {
				t.Skip("skipping AddTorrents: pass -torrent to test adding a torrent")
			}

			path := "/data/downloads"
			autoTMM := false
			err := client.AddTorrents(ctx, models.AddTorrentsRequest{
				Files:    []string{*flagTorrent},
				SavePath: &path,
				AutoTMM:  &autoTMM,
			})
			if err != nil {
				t.Fatalf("AddTorrents() error = %v", err)
			}
		})
	})
}
