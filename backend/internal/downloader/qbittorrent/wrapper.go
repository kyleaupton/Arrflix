package qbittorrent

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/kyleaupton/snaggle/backend/internal/downloader"
	qbt "github.com/superturkey650/go-qbittorrent/qbt"
)

const (
	maxRetries = 2
	retryDelay = 1 * time.Second
)

// qBittorrentClient wraps the library client to implement downloader.Client interface
type qBittorrentClient struct {
	instanceID downloader.InstanceID
	client     *qbt.Client
	username   string
	password   string
}

// NewQBittorrentClient creates a new qBittorrent client wrapper
func NewQBittorrentClient(instanceID downloader.InstanceID, baseURL, username, password string) *qBittorrentClient {
	return &qBittorrentClient{
		instanceID: instanceID,
		client:     qbt.NewClient(baseURL),
		username:   username,
		password:   password,
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

// Test tests the connection to qBittorrent
func (c *qBittorrentClient) Test(ctx context.Context) (downloader.TestResult, error) {
	result := downloader.TestResult{}

	// Test 1: Attempt to login
	if err := c.ensureLoggedIn(ctx); err != nil {
		// Determine error type
		errMsg := err.Error()
		var errorType string
		if strings.Contains(errMsg, "timeout") || strings.Contains(errMsg, "connection refused") || strings.Contains(errMsg, "no such host") {
			errorType = "Unable to connect to qBittorrent. Check if qBittorrent is running and the URL is correct."
		} else if strings.Contains(errMsg, "401") || strings.Contains(errMsg, "403") || strings.Contains(errMsg, "unauthorized") || strings.Contains(errMsg, "Fails") {
			errorType = "Authentication failed - check username and password"
		} else {
			errorType = "Connection test failed: " + errMsg
		}

		result.Success = false
		result.Error = errorType
		return result, nil // Return result with error, don't wrap in error
	}

	// Test 2: Get application version
	// Make direct HTTP call to check status code since library doesn't validate it
	version, err := c.getVersionWithStatusCheck(ctx)
	if err != nil {
		result.Success = false
		result.Error = "Connected but unable to retrieve version information: " + err.Error()
		return result, nil
	}

	// Test 3: Get Web API version (optional, don't fail if it errors)
	webAPIVersion, _ := c.getWebAPIVersionWithStatusCheck(ctx)

	// Success!
	result.Success = true
	result.Message = "Connection test successful"
	result.Version = version
	result.WebAPIVersion = webAPIVersion
	return result, nil
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

		// Snapshot existing torrents so we can identify a newly-added torrent even when
		// we can't extract a hash (e.g. when adding a .torrent URL).
		existing := map[string]bool{}
		if torrents, listErr := c.client.Torrents(qbt.TorrentsOptions{}); listErr == nil {
			for _, t := range torrents {
				existing[t.Hash] = true
			}
		}

		// Build download options
		opts := qbt.DownloadOptions{}
		if req.SavePath != "" {
			opts.Savepath = &req.SavePath
		}
		if req.Category != "" {
			opts.Category = &req.Category
		}
		if req.Paused {
			paused := true
			opts.Paused = &paused
		}

		// Add the torrent using library
		err = c.client.DownloadLinks([]string{torrentURL}, opts)
		if err != nil {
			// If it's an auth error, clear session and retry
			if strings.Contains(err.Error(), "login") || strings.Contains(err.Error(), "401") || strings.Contains(err.Error(), "403") || strings.Contains(err.Error(), "unauthorized") {
				c.client.Authenticated = false
				continue
			}
			continue
		}

		// Extract hash from magnet URL or derive it from the client after add.
		hash, err := extractHashFromMagnet(torrentURL)
		if err != nil {
			// If this isn't a magnet (e.g. .torrent URL), locate the newly-added torrent by diffing hashes.
			torrents, listErr := c.client.Torrents(qbt.TorrentsOptions{})
			if listErr != nil {
				return result, fmt.Errorf("failed to list torrents after add: %w", listErr)
			}
			var newest *qbt.TorrentInfo
			for i := range torrents {
				t := &torrents[i]
				if existing[t.Hash] {
					continue
				}
				// First unseen becomes candidate; if multiple, pick newest.
				if newest == nil || t.AddedOn > newest.AddedOn {
					newest = t
				}
			}
			if newest == nil {
				// Fall back to name-based lookup for magnets without btih (rare) or if diffing failed.
				hash, err = c.getHashFromName(ctx, req.MagnetURL)
				if err != nil {
					return result, fmt.Errorf("failed to get torrent hash: %w", err)
				}
			} else {
				hash = newest.Hash
				result.Name = newest.Name
			}
		}

		// Add tags if provided
		if len(req.Tags) > 0 {
			_, tagErr := c.client.AddTorrentTags([]string{hash}, req.Tags)
			if tagErr != nil {
				// Log but don't fail - tags are optional
				// Could log this if we had a logger
			}
		}

		result.ExternalID = hash
		if result.Name == "" {
			result.Name = extractNameFromMagnet(torrentURL)
		}

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

		// Use Torrents with hash filter
		opts := qbt.TorrentsOptions{
			Hashes: []string{externalID},
		}
		torrents, err := c.client.Torrents(opts)
		if err != nil {
			if strings.Contains(err.Error(), "login") || strings.Contains(err.Error(), "401") || strings.Contains(err.Error(), "403") || strings.Contains(err.Error(), "unauthorized") {
				c.client.Authenticated = false
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
		item.AddedAt = time.Unix(t.AddedOn, 0)

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

		// Get all torrents
		torrents, err := c.client.Torrents(qbt.TorrentsOptions{})
		if err != nil {
			if strings.Contains(err.Error(), "login") || strings.Contains(err.Error(), "401") || strings.Contains(err.Error(), "403") || strings.Contains(err.Error(), "unauthorized") {
				c.client.Authenticated = false
				continue
			}
			continue
		}

		items = make([]downloader.Item, len(torrents))
		for i, t := range torrents {
			items[i] = downloader.Item{
				ExternalID:  t.Hash,
				Name:        t.Name,
				Status:      mapStateToStatus(t.State),
				Progress:    t.Progress,
				SavePath:    t.SavePath,
				ContentPath: t.ContentPath,
				AddedAt:     time.Unix(t.AddedOn, 0),
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

		// Use library's TorrentFiles method
		torrentFiles, err := c.client.TorrentFiles(externalID)
		if err != nil {
			if strings.Contains(err.Error(), "login") || strings.Contains(err.Error(), "401") || strings.Contains(err.Error(), "403") || strings.Contains(err.Error(), "unauthorized") {
				c.client.Authenticated = false
				continue
			}
			continue
		}

		files = make([]downloader.File, len(torrentFiles))
		for i, f := range torrentFiles {
			files[i] = downloader.File{
				Path:     f.Name,
				Size:     int64(f.Size),
				Progress: f.Progress,
				Priority: f.Priority,
			}
		}

		return files, nil
	}

	return nil, fmt.Errorf("failed to list files after %d attempts: %w", maxRetries+1, err)
}

// Pause pauses a torrent
func (c *qBittorrentClient) Pause(ctx context.Context, externalID string) error {
	return c.withRetry(ctx, func() error {
		return c.client.Pause([]string{externalID})
	})
}

// Resume resumes a torrent
func (c *qBittorrentClient) Resume(ctx context.Context, externalID string) error {
	return c.withRetry(ctx, func() error {
		return c.client.Resume([]string{externalID})
	})
}

// Remove removes a torrent
func (c *qBittorrentClient) Remove(ctx context.Context, externalID string, deleteData bool) error {
	return c.withRetry(ctx, func() error {
		return c.client.Delete([]string{externalID}, deleteData)
	})
}

// ensureLoggedIn ensures the client is logged in
func (c *qBittorrentClient) ensureLoggedIn(ctx context.Context) error {
	if !c.client.Authenticated {
		return c.client.Login(c.username, c.password)
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

		if strings.Contains(err.Error(), "login") || strings.Contains(err.Error(), "401") || strings.Contains(err.Error(), "403") || strings.Contains(err.Error(), "unauthorized") {
			c.client.Authenticated = false
			continue
		}
	}
	return fmt.Errorf("failed after %d attempts: %w", maxRetries+1, err)
}

// getHashFromName tries to find a torrent hash by searching for the name
func (c *qBittorrentClient) getHashFromName(ctx context.Context, magnetURL string) (string, error) {
	name := extractNameFromMagnet(magnetURL)
	if name == "" {
		return "", fmt.Errorf("could not extract name from magnet URL")
	}

	torrents, err := c.client.Torrents(qbt.TorrentsOptions{})
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
	// Parse magnet URL
	if !strings.HasPrefix(magnetURL, "magnet:") {
		return "", fmt.Errorf("not a magnet URL")
	}

	// Extract hash from magnet URL format: magnet:?xt=urn:btih:HASH
	parts := strings.Split(magnetURL, "?")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid magnet URL format")
	}

	query := parts[1]
	params := strings.Split(query, "&")
	for _, param := range params {
		if strings.HasPrefix(param, "xt=urn:btih:") {
			hash := strings.TrimPrefix(param, "xt=urn:btih:")
			// Hash might be followed by & or end of string
			if idx := strings.Index(hash, "&"); idx != -1 {
				hash = hash[:idx]
			}
			if len(hash) != 40 && len(hash) != 32 {
				return "", fmt.Errorf("invalid hash length")
			}
			return hash, nil
		}
	}

	return "", fmt.Errorf("no hash found in magnet URL")
}

// extractNameFromMagnet extracts the name from a magnet URL
func extractNameFromMagnet(magnetURL string) string {
	// Parse magnet URL to extract dn parameter
	u, err := url.Parse(magnetURL)
	if err != nil {
		return ""
	}

	// Get dn parameter from query string
	name := u.Query().Get("dn")
	if name != "" {
		// URL decode the name
		decoded, err := url.QueryUnescape(name)
		if err == nil {
			return decoded
		}
		return name
	}
	return ""
}

// getVersionWithStatusCheck makes a direct HTTP request to check status code
// This works around the library's issue where it returns "Forbidden" as a string instead of an error
func (c *qBittorrentClient) getVersionWithStatusCheck(ctx context.Context) (string, error) {
	// Build the URL
	versionURL := strings.TrimSuffix(c.client.URL, "/") + "/api/v2/app/version"

	req, err := http.NewRequestWithContext(ctx, "GET", versionURL, nil)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	// Create HTTP client with same cookie jar to preserve authentication
	httpClient := &http.Client{
		Jar: c.client.Jar,
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check status code - this is the key fix
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("authentication failed: received status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	version := strings.TrimSpace(string(body))
	// Additional validation: version should not be empty and should look like a version
	if version == "" {
		return "", fmt.Errorf("empty version response")
	}

	return version, nil
}

// getWebAPIVersionWithStatusCheck makes a direct HTTP request to check status code
func (c *qBittorrentClient) getWebAPIVersionWithStatusCheck(ctx context.Context) (string, error) {
	// Build the URL
	versionURL := strings.TrimSuffix(c.client.URL, "/") + "/api/v2/app/webapiVersion"

	req, err := http.NewRequestWithContext(ctx, "GET", versionURL, nil)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	// Create HTTP client with same cookie jar to preserve authentication
	httpClient := &http.Client{
		Jar: c.client.Jar,
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("authentication failed: received status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	return strings.TrimSpace(string(body)), nil
}
