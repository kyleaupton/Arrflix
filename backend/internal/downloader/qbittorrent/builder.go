package qbittorrent

import (
	"encoding/json"
	"fmt"

	"github.com/kyleaupton/Arrflix/internal/downloader"
)

// Config represents qBittorrent-specific configuration
type Config struct {
	// Future: Add qBittorrent-specific config options here
	// For now, we use the base URL, username, and password from the DB record
}

// Build creates a qBittorrent client from a config record
func Build(rec downloader.ConfigRecord) (downloader.Client, error) {
	// Parse config JSON if present (for future qBittorrent-specific options)
	var config Config
	if len(rec.Config) > 0 {
		if err := json.Unmarshal(rec.Config, &config); err != nil {
			return nil, fmt.Errorf("parse config JSON: %w", err)
		}
	}

	// Extract credentials
	username := ""
	password := ""
	if rec.Username != nil {
		username = *rec.Username
	}
	if rec.Password != nil {
		password = *rec.Password
	}

	// Create the client wrapper
	client := NewQBittorrentClient(rec.ID, rec.URL, username, password)
	return client, nil
}

// Register registers the qBittorrent builder with the registry
func Register(registry *downloader.Registry) {
	registry.Register(downloader.TypeQbittorrent, Build)
}
