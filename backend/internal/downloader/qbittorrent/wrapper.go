package qbittorrent

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/kyleaupton/snaggle/backend/internal/downloader"
)

const (
	maxRetries = 2
	retryDelay = 1 * time.Second
)

// qBittorrentClient wraps the HTTP client to implement downloader.Client interface
type qBittorrentClient struct {
	instanceID downloader.InstanceID
	client     *Client
}

// NewQBittorrentClient creates a new qBittorrent client wrapper
func NewQBittorrentClient(instanceID downloader.InstanceID, baseURL, username, password string) *qBittorrentClient {
	return &qBittorrentClient{
		instanceID: instanceID,
		client:     NewClient(baseURL, username, password),
	}
}

// Type returns the downloader type
func (c *qBittorrentClient) Type() downloader.Type {
	return downloader.TypeQbittorrent
}

// InstanceID returns the instance ID
func (c *qBittorrentClient) InstanceID() downloader.InstanceID {
	return c.instanceID
}

// Add adds a download (magnet URL or torrent file)
func (c *qBittorrentClient) Add(ctx context.Context, req downloader.AddRequest) (downloader.AddResult, error) {
	var err error
	var result downloader.AddResult

	// Determine the URL to use
	torrentURL := req.MagnetURL
	if torrentURL == "" {
		return result, fmt.Errorf("magnet URL is required")
	}

	// Retry logic
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(retryDelay)
		}

		// Ensure we're logged in
		if err = c.ensureLoggedIn(ctx); err != nil {
			continue
		}

		// Add the torrent
		err = c.client.AddTorrent(ctx, torrentURL, req.SavePath, req.Category, req.Tags)
		if err != nil {
			// If it's an auth error, clear session and retry
			if strings.Contains(err.Error(), "login") || strings.Contains(err.Error(), "401") || strings.Contains(err.Error(), "403") {
				c.client.sid = ""
				continue
			}
			continue
		}

		// Extract hash from magnet URL or get it from the torrent
		hash, err := extractHashFromMagnet(torrentURL)
		if err != nil {
			// Try to get it from the torrent list
			hash, err = c.getHashFromName(ctx, req.MagnetURL)
			if err != nil {
				return result, fmt.Errorf("failed to get torrent hash: %w", err)
			}
		}

		result.ExternalID = hash
		result.Name = extractNameFromMagnet(torrentURL)

		return result, nil
	}

	return result, fmt.Errorf("failed to add torrent after %d attempts: %w", maxRetries+1, err)
}

// Get gets a torrent by external ID (hash)
func (c *qBittorrentClient) Get(ctx context.Context, externalID string) (downloader.Item, error) {
	var err error
	var item downloader.Item

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(retryDelay)
		}

		if err = c.ensureLoggedIn(ctx); err != nil {
			continue
		}

		torrents, err := c.getTorrents(ctx, externalID)
		if err != nil {
			if strings.Contains(err.Error(), "login") || strings.Contains(err.Error(), "401") || strings.Contains(err.Error(), "403") {
				c.client.sid = ""
				continue
			}
			continue
		}

		if len(torrents) == 0 {
			return item, fmt.Errorf("torrent not found: %s", externalID)
		}

		t := torrents[0]
		item.ExternalID = t.Hash
		item.Name = t.Name
		item.Status = mapStateToStatus(t.State)
		item.Progress = t.Progress
		item.SavePath = t.SavePath
		item.ContentPath = t.ContentPath
		item.AddedAt = t.AddedOn

		return item, nil
	}

	return item, fmt.Errorf("failed to get torrent after %d attempts: %w", maxRetries+1, err)
}

// List lists all items (torrents) in the client
func (c *qBittorrentClient) List(ctx context.Context) ([]downloader.Item, error) {
	var err error
	var items []downloader.Item

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(retryDelay)
		}

		if err = c.ensureLoggedIn(ctx); err != nil {
			continue
		}

		torrents, err := c.getTorrents(ctx, "")
		if err != nil {
			if strings.Contains(err.Error(), "login") || strings.Contains(err.Error(), "401") || strings.Contains(err.Error(), "403") {
				c.client.sid = ""
				continue
			}
			continue
		}

		items = make([]downloader.Item, len(torrents))
		for i, t := range torrents {
			items[i] = downloader.Item{
				ExternalID:  t.Hash,
				Name:         t.Name,
				Status:       mapStateToStatus(t.State),
				Progress:     t.Progress,
				SavePath:     t.SavePath,
				ContentPath:   t.ContentPath,
				AddedAt:      t.AddedOn,
			}
		}

		return items, nil
	}

	return nil, fmt.Errorf("failed to list torrents after %d attempts: %w", maxRetries+1, err)
}

// ListFiles lists files for a torrent
func (c *qBittorrentClient) ListFiles(ctx context.Context, externalID string) ([]downloader.File, error) {
	var err error
	var files []downloader.File

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(retryDelay)
		}

		if err = c.ensureLoggedIn(ctx); err != nil {
			continue
		}

		files, err = c.getTorrentFiles(ctx, externalID)
		if err != nil {
			if strings.Contains(err.Error(), "login") || strings.Contains(err.Error(), "401") || strings.Contains(err.Error(), "403") {
				c.client.sid = ""
				continue
			}
			continue
		}

		return files, nil
	}

	return nil, fmt.Errorf("failed to list files after %d attempts: %w", maxRetries+1, err)
}

// Pause pauses a torrent
func (c *qBittorrentClient) Pause(ctx context.Context, externalID string) error {
	return c.withRetry(ctx, func() error {
		return c.pauseTorrent(ctx, externalID)
	})
}

// Resume resumes a torrent
func (c *qBittorrentClient) Resume(ctx context.Context, externalID string) error {
	return c.withRetry(ctx, func() error {
		return c.resumeTorrent(ctx, externalID)
	})
}

// Remove removes a torrent
func (c *qBittorrentClient) Remove(ctx context.Context, externalID string, deleteData bool) error {
	return c.withRetry(ctx, func() error {
		return c.removeTorrent(ctx, externalID, deleteData)
	})
}

// ensureLoggedIn ensures the client is logged in
func (c *qBittorrentClient) ensureLoggedIn(ctx context.Context) error {
	if c.client.sid == "" {
		return c.client.Login(ctx)
	}
	return nil
}

// withRetry executes a function with retry logic
func (c *qBittorrentClient) withRetry(ctx context.Context, fn func() error) error {
	var err error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(retryDelay)
		}

		if err = c.ensureLoggedIn(ctx); err != nil {
			continue
		}

		err = fn()
		if err == nil {
			return nil
		}

		if strings.Contains(err.Error(), "login") || strings.Contains(err.Error(), "401") || strings.Contains(err.Error(), "403") {
			c.client.sid = ""
			continue
		}
	}
	return fmt.Errorf("failed after %d attempts: %w", maxRetries+1, err)
}

// getTorrents gets torrents from qBittorrent API
func (c *qBittorrentClient) getTorrents(ctx context.Context, hash string) ([]torrentInfo, error) {
	torrentsURL := fmt.Sprintf("%s/api/v2/torrents/info", c.client.baseURL)
	if hash != "" {
		torrentsURL += "?hashes=" + hash
	}

	req, err := http.NewRequestWithContext(ctx, "GET", torrentsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	c.client.addAuthHeader(req)

	resp, err := c.client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var torrents []torrentInfo
	if err := json.NewDecoder(resp.Body).Decode(&torrents); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return torrents, nil
}

// getTorrentFiles gets files for a torrent
func (c *qBittorrentClient) getTorrentFiles(ctx context.Context, hash string) ([]downloader.File, error) {
	filesURL := fmt.Sprintf("%s/api/v2/torrents/files?hash=%s", c.client.baseURL, hash)

	req, err := http.NewRequestWithContext(ctx, "GET", filesURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	c.client.addAuthHeader(req)

	resp, err := c.client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var files []torrentFile
	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	result := make([]downloader.File, len(files))
	for i, f := range files {
		result[i] = downloader.File{
			Path:     f.Name,
			Size:     f.Size,
			Progress: f.Progress,
			Priority: f.Priority,
		}
	}

	return result, nil
}

// pauseTorrent pauses a torrent
func (c *qBittorrentClient) pauseTorrent(ctx context.Context, hash string) error {
	return c.controlTorrent(ctx, "pause", hash)
}

// resumeTorrent resumes a torrent
func (c *qBittorrentClient) resumeTorrent(ctx context.Context, hash string) error {
	return c.controlTorrent(ctx, "resume", hash)
}

// removeTorrent removes a torrent
func (c *qBittorrentClient) removeTorrent(ctx context.Context, hash string, deleteData bool) error {
	action := "delete"
	if deleteData {
		action = "deletePerm"
	}
	return c.controlTorrent(ctx, action, hash)
}

// controlTorrent sends a control command to qBittorrent
func (c *qBittorrentClient) controlTorrent(ctx context.Context, action, hash string) error {
	controlURL := fmt.Sprintf("%s/api/v2/torrents/%s", c.client.baseURL, action)

	data := url.Values{}
	data.Set("hashes", hash)

	req, err := http.NewRequestWithContext(ctx, "POST", controlURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.client.addAuthHeader(req)

	resp, err := c.client.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, _ := io.ReadAll(resp.Body)
	responseText := strings.TrimSpace(string(body))
	if responseText != "Ok." {
		return fmt.Errorf("control failed: %s", responseText)
	}

	return nil
}

// getHashFromName tries to find a torrent hash by searching for the name
func (c *qBittorrentClient) getHashFromName(ctx context.Context, magnetURL string) (string, error) {
	name := extractNameFromMagnet(magnetURL)
	if name == "" {
		return "", fmt.Errorf("could not extract name from magnet URL")
	}

	torrents, err := c.getTorrents(ctx, "")
	if err != nil {
		return "", err
	}

	for _, t := range torrents {
		if t.Name == name {
			return t.Hash, nil
		}
	}

	return "", fmt.Errorf("torrent not found by name")
}

// mapStateToStatus maps qBittorrent state to JobStatus
func mapStateToStatus(state string) downloader.JobStatus {
	// qBittorrent states: downloading, seeding, completed, paused, queued, checking, error, missingFiles
	switch state {
	case "downloading":
		return downloader.StatusDownloading
	case "seeding":
		return downloader.StatusSeeding
	case "completed":
		return downloader.StatusCompleted
	case "pausedDL", "pausedUP":
		return downloader.StatusPaused
	case "queuedDL", "queuedUP":
		return downloader.StatusQueued
	case "error", "missingFiles":
		return downloader.StatusErrored
	default:
		return downloader.StatusUnknown
	}
}

// extractHashFromMagnet extracts the hash from a magnet URL
func extractHashFromMagnet(magnetURL string) (string, error) {
	u, err := url.Parse(magnetURL)
	if err != nil {
		return "", err
	}

	// Magnet URLs have format: magnet:?xt=urn:btih:HASH
	if u.Scheme != "magnet" {
		return "", fmt.Errorf("not a magnet URL")
	}

	xt := u.Query().Get("xt")
	if xt == "" {
		return "", fmt.Errorf("no xt parameter")
	}

	// Extract hash from urn:btih:HASH
	parts := strings.Split(xt, ":")
	if len(parts) < 3 {
		return "", fmt.Errorf("invalid xt parameter")
	}

	hash := parts[len(parts)-1]
	if len(hash) != 40 && len(hash) != 32 {
		return "", fmt.Errorf("invalid hash length")
	}

	return hash, nil
}

// extractNameFromMagnet extracts the name from a magnet URL
func extractNameFromMagnet(magnetURL string) string {
	u, err := url.Parse(magnetURL)
	if err != nil {
		return ""
	}

	return u.Query().Get("dn")
}

// torrentInfo represents a torrent from qBittorrent API
type torrentInfo struct {
	Hash         string    `json:"hash"`
	Name         string    `json:"name"`
	Size         int64     `json:"size"`
	Progress     float64   `json:"progress"`
	State        string    `json:"state"`
	SavePath     string    `json:"save_path"`
	ContentPath  string    `json:"content_path"`
	AddedOn      time.Time `json:"added_on"`
}

// torrentFile represents a file in a torrent from qBittorrent API
type torrentFile struct {
	Name     string  `json:"name"`
	Size     int64   `json:"size"`
	Progress float64 `json:"progress"`
	Priority int     `json:"priority"`
}

