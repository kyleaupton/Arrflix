package qbittorrent

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client is a qBittorrent Web API v2 client
type Client struct {
	baseURL    string
	username   string
	password   string
	httpClient *http.Client
	sid        string // session ID cookie
}

// NewClient creates a new qBittorrent client
func NewClient(baseURL, username, password string) *Client {
	return &Client{
		baseURL:    strings.TrimSuffix(baseURL, "/"),
		username:   username,
		password:   password,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// Login authenticates with qBittorrent and stores the session cookie
func (c *Client) Login(ctx context.Context) error {
	loginURL := fmt.Sprintf("%s/api/v2/auth/login", c.baseURL)

	data := url.Values{}
	data.Set("username", c.username)
	data.Set("password", c.password)

	req, err := http.NewRequestWithContext(ctx, "POST", loginURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("create login request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("login request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read login response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Check response body - qBittorrent returns "Ok." on success, "Fails." on failure
	responseText := strings.TrimSpace(string(body))
	if responseText != "Ok." {
		return fmt.Errorf("login failed: %s", responseText)
	}

	// Extract session cookie (SID)
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "SID" {
			c.sid = cookie.Value
			return nil
		}
	}

	// If no SID cookie, try to get it from Set-Cookie header
	setCookie := resp.Header.Get("Set-Cookie")
	if setCookie != "" {
		parts := strings.Split(setCookie, ";")
		for _, part := range parts {
			if strings.HasPrefix(strings.TrimSpace(part), "SID=") {
				c.sid = strings.TrimPrefix(strings.TrimSpace(part), "SID=")
				return nil
			}
		}
	}

	return fmt.Errorf("no session cookie received")
}

// Logout logs out from qBittorrent
func (c *Client) Logout(ctx context.Context) error {
	logoutURL := fmt.Sprintf("%s/api/v2/auth/logout", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, "POST", logoutURL, nil)
	if err != nil {
		return fmt.Errorf("create logout request: %w", err)
	}

	c.addAuthHeader(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("logout request failed: %w", err)
	}
	defer resp.Body.Close()

	c.sid = ""
	return nil
}

// AddTorrent adds a torrent to qBittorrent
func (c *Client) AddTorrent(ctx context.Context, torrentURL, savePath string, category string, tags []string) error {
	if c.sid == "" {
		if err := c.Login(ctx); err != nil {
			return fmt.Errorf("login required: %w", err)
		}
	}

	addURL := fmt.Sprintf("%s/api/v2/torrents/add", c.baseURL)

	data := url.Values{}
	data.Set("urls", torrentURL)
	if savePath != "" {
		data.Set("savepath", savePath)
	}
	if category != "" {
		data.Set("category", category)
	}
	if len(tags) > 0 {
		data.Set("tags", strings.Join(tags, ","))
	}

	req, err := http.NewRequestWithContext(ctx, "POST", addURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("create add torrent request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.addAuthHeader(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("add torrent request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read add torrent response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("add torrent failed with status %d: %s", resp.StatusCode, string(body))
	}

	responseText := strings.TrimSpace(string(body))
	if responseText != "Ok." {
		return fmt.Errorf("add torrent failed: %s", responseText)
	}

	return nil
}

// GetVersion returns the qBittorrent version
func (c *Client) GetVersion(ctx context.Context) (string, error) {
	if c.sid == "" {
		if err := c.Login(ctx); err != nil {
			return "", fmt.Errorf("login required: %w", err)
		}
	}

	versionURL := fmt.Sprintf("%s/api/v2/app/version", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", versionURL, nil)
	if err != nil {
		return "", fmt.Errorf("create version request: %w", err)
	}

	c.addAuthHeader(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("version request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read version response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("get version failed with status %d: %s", resp.StatusCode, string(body))
	}

	return strings.TrimSpace(string(body)), nil
}

// Torrent represents a torrent in qBittorrent
type Torrent struct {
	Hash     string
	Name     string
	Size     int64
	Progress float64
	State    string
}

// GetTorrents returns all torrents from qBittorrent
func (c *Client) GetTorrents(ctx context.Context) ([]Torrent, error) {
	if c.sid == "" {
		if err := c.Login(ctx); err != nil {
			return nil, fmt.Errorf("login required: %w", err)
		}
	}

	torrentsURL := fmt.Sprintf("%s/api/v2/torrents/info", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", torrentsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create get torrents request: %w", err)
	}

	c.addAuthHeader(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get torrents request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read get torrents response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get torrents failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse JSON response
	// For now, return empty slice - full parsing can be added later if needed
	_ = body
	return []Torrent{}, nil
}

// addAuthHeader adds the authentication cookie to the request
func (c *Client) addAuthHeader(req *http.Request) {
	if c.sid != "" {
		req.AddCookie(&http.Cookie{
			Name:  "SID",
			Value: c.sid,
		})
	}
}
