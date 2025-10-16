package jackett

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"runtime/debug"
	"strings"
	"sync"

	"golang.org/x/net/context"
)

func getXML[T any](ctx context.Context, client *http.Client, url string) (T, error) {
	var data T
	b, err := getBytes(ctx, client, url)
	if err != nil {
		return data, err
	}
	if err := xml.Unmarshal(b, &data); err != nil {
		return data, fmt.Errorf("unmarshal response data: %s: %w\n%s", url, err, string(b))
	}
	return data, nil
}

func getJSON[T any](ctx context.Context, client *http.Client, url string) (T, error) {
	var data T
	b, err := getBytes(ctx, client, url)

	test, _ := json.MarshalIndent(b, "", "  ")
	fmt.Printf(string(test))

	if err != nil {
		return data, err
	}

	// Log the raw response for debugging
	fmt.Printf("DEBUG: Raw response from %s:\n", url)
	fmt.Printf("Response length: %d bytes\n", len(b))
	fmt.Printf("First 500 chars: %s\n", string(b[:int(math.Min(500, float64(len(b))))]))

	if err := json.Unmarshal(b, &data); err != nil {
		return data, fmt.Errorf("unmarshal response data: %s: %w\n%s", url, err, string(b))
	}
	return data, nil
}

func getBytes(ctx context.Context, client *http.Client, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("make fetch request: %w", err)
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("invoke fetch request: %w", err)
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read respose: %w", err)
	}
	return b, nil
}

func postJSON[T any](ctx context.Context, client *http.Client, url string, data interface{}) (T, error) {
	var result T
	jsonData, err := json.Marshal(data)
	if err != nil {
		return result, fmt.Errorf("marshal request data: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return result, fmt.Errorf("make POST request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return result, fmt.Errorf("invoke POST request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return result, fmt.Errorf("read response: %w", err)
	}

	if res.StatusCode >= 400 {
		return result, fmt.Errorf("HTTP %d: %s", res.StatusCode, string(body))
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return result, fmt.Errorf("unmarshal response data: %s: %w\n%s", url, err, string(body))
	}
	return result, nil
}

func putJSON[T any](ctx context.Context, client *http.Client, url string, data interface{}) (T, error) {
	var result T
	jsonData, err := json.Marshal(data)
	if err != nil {
		return result, fmt.Errorf("marshal request data: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return result, fmt.Errorf("make PUT request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return result, fmt.Errorf("invoke PUT request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return result, fmt.Errorf("read response: %w", err)
	}

	if res.StatusCode >= 400 {
		return result, fmt.Errorf("HTTP %d: %s", res.StatusCode, string(body))
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return result, fmt.Errorf("unmarshal response data: %s: %w\n%s", url, err, string(body))
	}
	return result, nil
}

func deleteRequest(ctx context.Context, client *http.Client, url string) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("make DELETE request: %w", err)
	}

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("invoke DELETE request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("HTTP %d: %s", res.StatusCode, string(body))
	}

	return nil
}

var _ http.RoundTripper = (*middleware)(nil)

// wrapTransport wraps the given http.Transport with a middleware that adds the user agent to all outgoing requests. It also adds the api key and/or cookies to all requests matching BaseURL.
func wrapTransport(rt http.RoundTripper, base *url.URL, apiKey string, client *Client) http.RoundTripper {
	if rt == nil {
		rt = http.DefaultTransport
	}
	return &middleware{
		Transport: rt,
		BaseURL:   base,
		APIKey:    apiKey,
		Client:    client,
	}
}

type middleware struct {
	Transport http.RoundTripper
	BaseURL   *url.URL
	APIKey    string
	Client    *Client
}

func (m *middleware) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("User-Agent", ua())

	if m.matchesTarget(r.URL) {
		// Add API key if provided
		if m.APIKey != "" {
			q := r.URL.Query()
			q.Set("apikey", m.APIKey)
			r.URL.RawQuery = q.Encode()
			fmt.Printf("DEBUG: Making request to %s with API key: %s\n", r.URL.String(), m.APIKey)
		}

		// Add cookies if provided
		if m.Client != nil {
			m.Client.cookieMutex.RLock()
			if len(m.Client.cfg.Cookies) > 0 {
				var cookieStrings []string
				for name, value := range m.Client.cfg.Cookies {
					cookieStrings = append(cookieStrings, fmt.Sprintf("%s=%s", name, value))
				}
				r.Header.Set("Cookie", strings.Join(cookieStrings, "; "))
				fmt.Printf("DEBUG: Making request to %s with cookies: %s\n", r.URL.String(), strings.Join(cookieStrings, "; "))
			}
			m.Client.cookieMutex.RUnlock()
		}
	}

	// Make the request
	resp, err := m.Transport.RoundTrip(r)
	if err != nil {
		return resp, err
	}

	// Check for authentication failures and retry with fresh cookies
	if m.matchesTarget(r.URL) && m.Client != nil && isAuthFailure(resp.StatusCode) {
		// Check if this is a retry (to prevent infinite loops)
		if r.Header.Get("X-Retry-Attempt") == "1" {
			return resp, err
		}

		// Try to acquire fresh cookies
		if err := m.Client.acquireCookies(r.Context()); err != nil {
			fmt.Printf("DEBUG: Failed to acquire cookies: %v\n", err)
			return resp, err
		}

		// Retry the request with fresh cookies
		retryReq := r.Clone(r.Context())
		retryReq.Header.Set("X-Retry-Attempt", "1")

		// Add fresh cookies
		m.Client.cookieMutex.RLock()
		if len(m.Client.cfg.Cookies) > 0 {
			var cookieStrings []string
			for name, value := range m.Client.cfg.Cookies {
				cookieStrings = append(cookieStrings, fmt.Sprintf("%s=%s", name, value))
			}
			retryReq.Header.Set("Cookie", strings.Join(cookieStrings, "; "))
			fmt.Printf("DEBUG: Retrying request to %s with fresh cookies: %s\n", retryReq.URL.String(), strings.Join(cookieStrings, "; "))
		}
		m.Client.cookieMutex.RUnlock()

		// Close the original response body
		resp.Body.Close()

		// Make the retry request
		return m.Transport.RoundTrip(retryReq)
	}

	return resp, err
}

// isAuthFailure checks if the status code indicates an authentication failure
func isAuthFailure(statusCode int) bool {
	return statusCode == 302 || statusCode == 401 || statusCode == 403
}

func (m *middleware) matchesTarget(reqURL *url.URL) bool {
	if m.BaseURL == nil {
		return false
	}

	if reqURL.Scheme != m.BaseURL.Scheme {
		return false
	}

	return normalizeHost(reqURL) == normalizeHost(m.BaseURL)
}

func normalizeHost(u *url.URL) string {
	host := u.Host
	if !strings.Contains(host, ":") {
		switch u.Scheme {
		case "http":
			host += ":80"
		case "https":
			host += ":443"
		}
	}
	return strings.ToLower(host)
}

var ua = sync.OnceValue(func() string {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return "go-jackett/unknown"
	}
	version := buildInfo.Main.Version
	if version == "" || version == "(devel)" {
		version = "dev"
	}
	return fmt.Sprintf("go-jackett/%s", version)
})
