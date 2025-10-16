package jackett

import (
	"fmt"
	"math"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
)

// Settings defines the configuration for the Jackett client.
type Settings struct {
	// ApiURL is the base URL for the Jackett API.
	// If empty, the value of the JACKETT_API_URL environment variable will be used.
	ApiURL string

	// ApiKey is the API key for accessing the Jackett API.
	// If empty, the value of the JACKETT_API_KEY environment variable will be used.
	ApiKey string

	// Cookies is a map of cookie names to values for cookie-based authentication.
	// If empty, the value of the JACKETT_COOKIES environment variable will be used.
	// The environment variable should be in the format "name1=value1; name2=value2".
	Cookies map[string]string

	// Client is the HTTP client to use for making requests.
	// If nil, http.DefaultClient will be used.
	Client *http.Client

	// DefaultTrackers is a list of tracker IDs to use if a FetchRequest does not specify any.
	// If empty and a FetchRequest does not specify trackers, "all" trackers will be used.
	DefaultTrackers []string
}

// Client is a Jackett API client.
type Client struct {
	apiURL      *url.URL
	cfg         Settings
	tcache      sync.Map
	cookieMutex sync.RWMutex // Protects cfg.Cookies access
}

const (
	envAPIURL  = "JACKETT_API_URL"
	envAPIKey  = "JACKETT_API_KEY"
	envCookies = "JACKETT_COOKIES"
)

// New creates a new Jackett client with the given settings.
// It will return an error if the API URL cannot be parsed.
// Environment variables JACKETT_API_URL, JACKETT_API_KEY, and JACKETT_COOKIES can be used
// as fallbacks if ApiURL, ApiKey, or Cookies are not provided in Settings.
func New(s Settings) (*Client, error) {
	j := Client{cfg: s}
	apiURLStr := valOrEnv(s.ApiURL, envAPIURL)
	apiURL, err := url.Parse(apiURLStr)
	if err != nil {
		return nil, fmt.Errorf("parse url: %s: %w", apiURLStr, err)
	}
	j.apiURL = apiURL

	// Handle cookies from environment variable if not provided in settings
	if len(j.cfg.Cookies) == 0 {
		envCookiesValue := os.Getenv(envCookies)
		if envCookiesValue != "" {
			j.cfg.Cookies = parseCookiesFromEnv(envCookiesValue)
		}
	}

	if j.cfg.Client == nil {
		j.cfg.Client = http.DefaultClient
	}
	j.cfg.Client.Transport = wrapTransport(j.cfg.Client.Transport,
		j.apiURL,
		valOrEnv(s.ApiKey, envAPIKey),
		&j)
	return &j, nil
}

type fetchOpts struct {
	MaxConcurrency     int
	ProgressReportFunc func(complete uint, total uint)
}

// FetchOption is a function that configures a Fetch operation.
type FetchOption func(*fetchOpts)

// WithMaxConcurrency sets the maximum number of concurrent requests for a Fetch operation.
// If n is less than 1, it defaults to runtime.NumCPU().
func WithMaxConcurrency(n uint) FetchOption {
	return func(o *fetchOpts) {
		o.MaxConcurrency = int(n)
	}
}

// WithProgressReportFunc sets a callback function to report progress during a Fetch operation.
// The callback receives the number of completed requests and the total number of requests.
func WithProgressReportFunc(f func(complete uint, total uint)) FetchOption {
	return func(o *fetchOpts) {
		o.ProgressReportFunc = f
	}
}

// Fetch executes a fetch request against the Jackett API.
// It returns a slice containing the combined results of the fetch, or an error.
// Fetch will concurrently request data from multiple trackers when possible.
// Error type will be *UnsupportedError when the query is not supported by the tracker.
func (j *Client) Fetch(ctx context.Context, fr *FetchRequest, opts ...FetchOption) ([]Result, error) {
	var o fetchOpts
	for _, f := range opts {
		f(&o)
	}
	if o.MaxConcurrency < 1 {
		o.MaxConcurrency = runtime.NumCPU()
	}
	if o.ProgressReportFunc == nil {
		o.ProgressReportFunc = func(_, _ uint) {}
	}
	o.ProgressReportFunc(0, math.MaxInt32)

	urls, err := j.generateFetchURLs(fr)
	if err != nil {
		return nil, fmt.Errorf("generate urls: %w", err)
	}

	var wg errgroup.Group
	wg.SetLimit(o.MaxConcurrency)

	// Ensure that all selected trackers support this query
	for _, u := range urls {
		wg.Go(func() error {
			tracker := extractTracker(u.Path)
			if tracker == "all" {
				return nil // can't check caps of meta indexers
			}
			caps, err := j.getIndexerCaps(ctx, tracker)
			if err != nil {
				return fmt.Errorf("get indexer caps: %s: %w", tracker, err)
			}
			if err := caps.Validate(u.Query()); err != nil {
				return fmt.Errorf("%s does not support this query: %w", tracker, err)
			}
			return nil
		})
	}

	if err := wg.Wait(); err != nil {
		return nil, fmt.Errorf("at least one query was invalid: %w", err)
	}

	total := uint(len(urls))
	var complete atomic.Uint32

	o.ProgressReportFunc(0, total)
	results := make([][]Result, len(urls))
	for i, url := range urls {
		wg.Go(func() error {
			defer func() {
				o.ProgressReportFunc(uint(complete.Add(1)), total)
			}()
			resp, err := getXML[searchResp](ctx, j.cfg.Client, url.String())
			if err != nil {
				return fmt.Errorf("fetch: %s: %w", url.String(), err)
			}
			results[i], err = resp.Unmarshal()
			return err
		})
	}

	err = wg.Wait()
	resp := slices.Concat(results...)
	o.ProgressReportFunc(total, total)
	return resp, err
}

// ListIndexers returns a slice of all indexers on this Jackett instance.
func (j *Client) ListIndexers(ctx context.Context, configured *bool) ([]IndexerDetails, error) {
	u := *j.apiURL
	u.Path = "/api/v2.0/indexers/all/results/torznab"
	q := u.Query()
	q.Add("t", "indexers")
	if configured != nil {
		q.Add("configured", strconv.FormatBool(*configured))
	}
	u.RawQuery = q.Encode()
	idxs, err := getXML[indexersResp](ctx, j.cfg.Client, u.String())
	if err != nil {
		return nil, fmt.Errorf("list indexers: %w", err)
	}
	slices.SortFunc(idxs.Indexers, func(a, b IndexerDetails) int {
		return strings.Compare(a.ID, b.ID)
	})
	return idxs.Indexers, err
}

// GetIndexerConfig retrieves the configuration for a specific indexer.
func (j *Client) GetIndexerConfig(ctx context.Context, indexerID string) (*IndexerConfigResponse, error) {
	u := *j.apiURL
	u.Path = fmt.Sprintf("/api/v2.0/indexers/%s/config", indexerID)

	config, err := getJSON[IndexerConfigResponse](ctx, j.cfg.Client, u.String())
	if err != nil {
		return nil, fmt.Errorf("get indexer config: %w", err)
	}

	// if config.Error != "" {
	// 	return nil, fmt.Errorf("indexer config error: %s", config.Error)
	// }

	return &config, nil
}

// SaveIndexerConfig creates or updates an indexer configuration.
func (j *Client) SaveIndexerConfig(ctx context.Context, indexerID string, config any) error {
	u := *j.apiURL
	u.Path = fmt.Sprintf("/api/v2.0/indexers/%s/config", indexerID)

	_, err := putJSON[any](ctx, j.cfg.Client, u.String(), config)
	if err != nil {
		return fmt.Errorf("save indexer config: %w", err)
	}

	return nil
}

// DeleteIndexer removes an indexer by its ID.
func (j *Client) DeleteIndexer(ctx context.Context, indexerID string) error {
	u := *j.apiURL
	u.Path = fmt.Sprintf("/api/v2.0/indexers/%s", indexerID)

	if err := deleteRequest(ctx, j.cfg.Client, u.String()); err != nil {
		return fmt.Errorf("delete indexer: %w", err)
	}

	return nil
}

func (j *Client) getIndexerCaps(ctx context.Context, id string) (IndexerCaps, error) {
	if v, ok := j.tcache.Load(id); ok {
		return v.(IndexerCaps), nil
	}
	u := *j.apiURL
	u.Path = fmt.Sprintf("/api/v2.0/indexers/%s/results/torznab", id)
	q := u.Query()
	q.Add("t", "caps")
	u.RawQuery = q.Encode()
	caps, err := getXML[IndexerCaps](ctx, j.cfg.Client, u.String())
	if err != nil {
		return caps, fmt.Errorf("list indexers: %w", err)
	}
	j.tcache.Store(id, caps)
	return caps, nil
}

func (j *Client) generateFetchURLs(fr *FetchRequest) ([]url.URL, error) {
	trackers := fr.Trackers()
	if len(trackers) == 0 {
		trackers = j.cfg.DefaultTrackers
	}
	if len(trackers) == 0 {
		trackers = []string{"all"} // meta tracker
	}

	var urls []url.URL
	for _, tracker := range trackers {
		u := *j.apiURL
		u.Path = fmt.Sprintf("/api/v2.0/indexers/%s/results/torznab", tracker)
		q, err := fr.Values()
		if err != nil {
			return nil, fmt.Errorf("marshal url: %w", err)
		}
		u.RawQuery = q.Encode()
		urls = append(urls, u)
	}
	return urls, nil
}

func valOrEnv(v, env string) string {
	if v != "" {
		return v
	}
	return os.Getenv(env)
}

// parseCookiesFromEnv parses cookies from environment variable format "name1=value1; name2=value2"
func parseCookiesFromEnv(envValue string) map[string]string {
	cookies := make(map[string]string)
	if envValue == "" {
		return cookies
	}

	// Split by semicolon and parse each cookie
	cookiePairs := strings.Split(envValue, ";")
	for _, pair := range cookiePairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}

		// Split by first equals sign
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			cookies[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	return cookies
}

func extractTracker(path string) string {
	_, tracker, ok := strings.Cut(path, "api/v2.0/indexers/")
	if !ok {
		return ""
	}
	return strings.TrimSuffix(tracker, "/results/torznab")
}

// acquireCookies performs the Jackett cookie authentication flow
// by following the redirect chain and extracting cookies from Set-Cookie headers
func (j *Client) acquireCookies(ctx context.Context) error {
	j.cookieMutex.Lock()
	defer j.cookieMutex.Unlock()

	// Create a temporary client that doesn't follow redirects automatically
	// Use a clean transport to avoid middleware interference
	tempClient := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // Don't follow redirects automatically
		},
		Transport: http.DefaultTransport,
	}

	// Start the flow at /UI/Dashboard
	dashboardURL := *j.apiURL
	dashboardURL.Path = "/UI/Dashboard"

	// Create a cookie jar to collect cookies during the flow
	jar, err := cookiejar.New(nil)
	if err != nil {
		return fmt.Errorf("create cookie jar: %w", err)
	}
	tempClient.Jar = jar

	// Follow the redirect chain
	currentURL := dashboardURL.String()
	maxRedirects := 10 // Prevent infinite loops
	redirectCount := 0

	for redirectCount < maxRedirects {
		req, err := http.NewRequestWithContext(ctx, "GET", currentURL, nil)
		if err != nil {
			return fmt.Errorf("create request: %w", err)
		}

		resp, err := tempClient.Do(req)
		if err != nil {
			return fmt.Errorf("execute request: %w", err)
		}
		resp.Body.Close()

		// Check if we got a successful response (200)
		if resp.StatusCode == 200 {
			break
		}

		// Check if we got a redirect
		if resp.StatusCode >= 300 && resp.StatusCode < 400 {
			location := resp.Header.Get("Location")
			if location == "" {
				return fmt.Errorf("redirect without Location header")
			}

			// Parse the location URL
			redirectURL, err := url.Parse(location)
			if err != nil {
				return fmt.Errorf("parse redirect URL: %w", err)
			}

			// Make absolute URL if needed
			if !redirectURL.IsAbs() {
				baseURL, err := url.Parse(currentURL)
				if err != nil {
					return fmt.Errorf("parse base URL: %w", err)
				}
				redirectURL = baseURL.ResolveReference(redirectURL)
			}

			currentURL = redirectURL.String()
			redirectCount++
			continue
		}

		// If we get here, it's not a redirect and not success
		return fmt.Errorf("unexpected status code %d during cookie acquisition", resp.StatusCode)
	}

	if redirectCount >= maxRedirects {
		return fmt.Errorf("too many redirects during cookie acquisition")
	}

	// Extract cookies from the jar and store them in our config
	if j.cfg.Cookies == nil {
		j.cfg.Cookies = make(map[string]string)
	}

	// Get cookies for our API URL
	cookies := jar.Cookies(j.apiURL)
	for _, cookie := range cookies {
		j.cfg.Cookies[cookie.Name] = cookie.Value
	}

	// Also check for any cookies that might be set for the base domain
	baseURL := *j.apiURL
	baseURL.Path = ""
	baseURL.RawQuery = ""
	baseURL.Fragment = ""
	baseCookies := jar.Cookies(&baseURL)
	for _, cookie := range baseCookies {
		j.cfg.Cookies[cookie.Name] = cookie.Value
	}

	return nil
}
