package qbittorrent

const apiBase = "/api/v2"

func buildPath(path string) string {
	return apiBase + path
}
