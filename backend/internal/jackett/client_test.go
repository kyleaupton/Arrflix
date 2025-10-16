package jackett

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

const testAPIKey string = "abracadabra"
const testCookieValue string = "CfDJ8Der2JDZHqxIpwmLd8ZiKBfPcdLxd2ZjZGlin34qAKJfs4OSWdX-qqScYz-fMbWZRB3yyM4XmoLiIbd898EM5FewjQxid3Xw-7T-0pS37mlQ3S-UUlM27AWyRVy8W-JiFLVFTPxLF6MKnKoZ6CEbNrubTnN5K8-j5p5eNeOnJAsgfjtQ-8GpbvCLr0hIy0bDXCgfRFNZrenfsSJ0pOJup_QYDuYv0bmDr36pTBYnYxDKh6Uh_unnstHxYj9fHE6J0HIAs67srQo5_3MukBnClj4vkjuX21HpXwxs6UI8IGrw5gLYZnXJ0_-z-302UNdi3xI0jLqDu8Izs1DbVccLkNT0"

func TestGenerateURL(t *testing.T) {
	t.Parallel()

	j, srv := mockServer(t)
	defer srv.Close()

	tests := []struct {
		input *FetchRequest
		want  []string
	}{
		{NewRawSearch().Build(),
			[]string{"/api/v2.0/indexers/all/results/torznab?extended=1&t=search"}},
		{NewRawSearch().
			WithQuery("qqq").
			WithTrackers("aaa", "bbb").
			WithCategories(1, 2).
			Build(),
			[]string{
				"/api/v2.0/indexers/aaa/results/torznab?cat=1%2C2&extended=1&q=qqq&t=search",
				"/api/v2.0/indexers/bbb/results/torznab?cat=1%2C2&extended=1&q=qqq&t=search"}},
	}
	for _, test := range tests {
		got, err := j.generateFetchURLs(test.input)
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != len(test.want) {
			t.Fatalf("expected %d urls, got %d", len(test.want), len(got))
		}
		for i := range test.want {
			if !strings.HasSuffix(got[i].String(), test.want[i]) {
				t.Errorf("strings.HasSuffix(generateURL(%+v), %q), want %q", test.input, got[i].String(), test.want[i])
			}
		}
	}
}

func TestDefaultTrackers(t *testing.T) {
	t.Parallel()

	j, srv := mockServer(t)
	defer srv.Close()

	// No default trackers and none defined on query
	urls, err := j.generateFetchURLs(NewRawSearch().Build())
	if err != nil {
		t.Fatal(err)
	}
	if len(urls) != 1 {
		t.Fatalf("expected 1 url got %d", len(urls))
	}
	if !strings.Contains(urls[0].Path, "/all/") {
		t.Errorf("expected url to contain all tracker: %s", urls[0].Path)
	}

	// Default trackers but none defined on query
	j.cfg.DefaultTrackers = []string{"foo", "bar"}
	urls, err = j.generateFetchURLs(NewRawSearch().Build())
	if err != nil {
		t.Fatal(err)
	}
	if len(urls) != 2 {
		t.Fatalf("expected 2 urls got %d", len(urls))
	}
	if !strings.Contains(urls[0].Path, "/foo/") {
		t.Errorf("expected url to contain foo tracker: %s", urls[0].Path)
	}
	if !strings.Contains(urls[1].Path, "/bar/") {
		t.Errorf("expected url to contain bar tracker: %s", urls[1].Path)
	}

	// Default trackers but query overrides
	j.cfg.DefaultTrackers = []string{"foo", "bar"}
	urls, err = j.generateFetchURLs(NewRawSearch().
		WithTrackers("baz").
		Build())
	if err != nil {
		t.Fatal(err)
	}
	if len(urls) != 1 {
		t.Fatalf("expected 1 url got %d", len(urls))
	}
	if !strings.Contains(urls[0].Path, "/baz/") {
		t.Errorf("expected url to contain baz tracker: %s", urls[0].Path)
	}
}

func TestFetch(t *testing.T) {
	t.Parallel()

	j, srv := mockServer(t)
	defer srv.Close()

	input := NewRawSearch().Build()
	got, err := j.Fetch(t.Context(), input)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 4 {
		t.Errorf("len(Fetch(%+v).Results) = %v, want 4", input, len(got))
	}
}

func TestFetchProgress(t *testing.T) {
	t.Parallel()

	j, srv := mockServer(t)
	defer srv.Close()

	var trackers []string
	for i := range 100 {
		trackers = append(trackers, fmt.Sprintf("tracker-%d", i))
	}

	var gotComplete, gotTotal atomic.Uint32
	input := NewRawSearch().WithTrackers(trackers...).Build()
	got, err := j.Fetch(t.Context(), input, WithProgressReportFunc(func(complete, total uint) {
		gotComplete.Store(uint32(complete))
		gotTotal.Store(uint32(total))
	}))
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 100 {
		t.Errorf("len(Fetch(%+v).Results) = %v, want 100", input, len(got))
	}
	if gotComplete.Load() != gotTotal.Load() {
		t.Errorf("expected complete and total counts to be equal: %d != %d", gotComplete.Load(), gotTotal.Load())
	}
}

func TestListIndexers(t *testing.T) {
	t.Parallel()

	j, srv := mockServer(t)
	defer srv.Close()

	got, err := j.ListIndexers(t.Context(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d indexers, expected 1", len(got))
	}
	if got[0].ID != "my-indexer" {
		t.Errorf("expected indexer to have id %q; got %q", "my-indexer", got[0].ID)
	}
}

func mockServer(t *testing.T) (*Client, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v2.0/indexers/all/results/torznab":
			if r.URL.Query().Get("t") == "indexers" {
				w.Write([]byte(exampleIndexerResponse))
				return
			}
			w.Write([]byte(exampleResponse))
		default:
			if r.URL.Query().Get("t") == "caps" {
				w.Write([]byte(exampleCapsResponse))
				return
			}
			tracker := extractTracker(r.URL.Path)
			if tracker == "" {
				// Not a known request pattern
				w.WriteHeader(http.StatusNotFound)
				return
			}
			if tracker == "all" {
				t.Fatal("tracker 'all' should have been caught by an earlier case")
			}
			responseXML := fmt.Sprintf(exampleResponseTemplate, tracker)
			w.Header().Set("Content-Type", "application/xml")
			w.Write([]byte(responseXML))
			return
		}
	}))

	j, err := New(Settings{
		ApiURL: srv.URL,
		ApiKey: testAPIKey,
		Client: srv.Client(),
	})
	if err != nil {
		t.Fatal(err)
	}
	return j, srv
}

const exampleResponse = `
<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom" xmlns:torznab="http://torznab.com/schemas/2015/feed">
  <channel>
    <atom:link href="http://localhost:9117/" rel="self" type="application/rss+xml" />
    <title>Generic-Tracker (API)</title>
    <description>Generic tracker description</description>
    <link>https://generic-tracker.example/</link>
    <language>en-US</language>
    <category>search</category>
    <item>
      <title>Media S15E12 720p DSNP WEB-DL DDP 5.1 H.264-FLUX</title>
      <guid>https://generic-tracker.example/torrents/media-s15e12-720p-dsnp-web-dl-ddp-51-h264-flux.444921</guid>
      <jackettindexer id="generic-tracker-api">Generic-Tracker (API)</jackettindexer>
      <type>private</type>
      <comments>https://generic-tracker.example/torrents/media-s15e12-720p-dsnp-web-dl-ddp-51-h264-flux.444921</comments>
      <pubDate>Fri, 30 May 2025 22:50:38 +0000</pubDate>
      <size>398928401</size>
      <grabs>22</grabs>
      <description />
      <link>http://localhost:9117/dl/generic-tracker-api/?jackett_apikey=PLACEHOLDER_API_KEY&amp;path=PLACEHOLDER_ENCODED_PATH&amp;file=Media+S15E12+720p+DSNP+WEB-DL+DDP+5.1+H.264-FLUX</link>
      <category>5000</category>
      <category>100002</category>
      <enclosure url="http://localhost:9117/dl/generic-tracker-api/?jackett_apikey=PLACEHOLDER_API_KEY&amp;path=PLACEHOLDER_ENCODED_PATH&amp;file=Media+S15E12+720p+DSNP+WEB-DL+DDP+5.1+H.264-FLUX" length="398928401" type="application/x-bittorrent" />
      <torznab:attr name="category" value="5000" />
      <torznab:attr name="category" value="100002" />
      <torznab:attr name="imdb" value="1561755" />
      <torznab:attr name="imdbid" value="tt1561755" />
      <torznab:attr name="tmdbid" value="32726" />
      <torznab:attr name="seeders" value="20" />
      <torznab:attr name="peers" value="21" />
      <torznab:attr name="infohash" value="PLACEHOLDER_INFOHASH_1" />
      <torznab:attr name="downloadvolumefactor" value="1" />
      <torznab:attr name="uploadvolumefactor" value="1" />
    </item>
    <item>
      <title>Media S15E12 1080p DSNP WEB-DL DDP 5.1 H.264-FLUX</title>
      <guid>https://generic-tracker.example/torrents/media-s15e12-1080p-dsnp-web-dl-ddp-51-h264-flux.444918</guid>
      <jackettindexer id="generic-tracker-api">Generic-Tracker (API)</jackettindexer>
      <type>private</type>
      <comments>https://generic-tracker.example/torrents/media-s15e12-1080p-dsnp-web-dl-ddp-51-h264-flux.444918</comments>
      <pubDate>Fri, 30 May 2025 22:50:24 +0000</pubDate>
      <size>620177942</size>
      <grabs>61</grabs>
      <description />
      <link>http://localhost:9117/dl/generic-tracker-api/?jackett_apikey=PLACEHOLDER_API_KEY&amp;path=PLACEHOLDER_ENCODED_PATH&amp;file=Media+S15E12+1080p+DSNP+WEB-DL+DDP+5.1+H.264-FLUX</link>
      <category>5000</category>
      <category>100002</category>
      <enclosure url="http://localhost:9117/dl/generic-tracker-api/?jackett_apikey=PLACEHOLDER_API_KEY&amp;path=PLACEHOLDER_ENCODED_PATH&amp;file=Media+S15E12+1080p+DSNP+WEB-DL+DDP+5.1+H.264-FLUX" length="620177942" type="application/x-bittorrent" />
      <torznab:attr name="category" value="5000" />
      <torznab:attr name="category" value="100002" />
      <torznab:attr name="imdb" value="1561755" />
      <torznab:attr name="imdbid" value="tt1561755" />
      <torznab:attr name="tmdbid" value="32726" />
      <torznab:attr name="seeders" value="59" />
      <torznab:attr name="peers" value="60" />
      <torznab:attr name="infohash" value="PLACEHOLDER_INFOHASH_2" />
      <torznab:attr name="downloadvolumefactor" value="1" />
      <torznab:attr name="uploadvolumefactor" value="1" />
    </item>
    <item>
      <title>Media S15E12 1080p DSNP WEB-DL DDP 5.1 H.264-NTb</title>
      <guid>https://generic-tracker.example/torrents/media-s15e12-1080p-dsnp-web-dl-ddp-51-h264-ntb.444841</guid>
      <jackettindexer id="generic-tracker-api">Generic-Tracker (API)</jackettindexer>
      <type>private</type>
      <comments>https://generic-tracker.example/torrents/media-s15e12-1080p-dsnp-web-dl-ddp-51-h264-ntb.444841</comments>
      <pubDate>Fri, 30 May 2025 14:25:03 +0000</pubDate>
      <size>558112511</size>
      <grabs>97</grabs>
      <description />
      <link>http://localhost:9117/dl/generic-tracker-api/?jackett_apikey=PLACEHOLDER_API_KEY&amp;path=PLACEHOLDER_ENCODED_PATH&amp;file=Media+S15E12+1080p+DSNP+WEB-DL+DDP+5.1+H.264-NTb</link>
      <category>5000</category>
      <category>100002</category>
      <enclosure url="http://localhost:9117/dl/generic-tracker-api/?jackett_apikey=PLACEHOLDER_API_KEY&amp;path=PLACEHOLDER_ENCODED_PATH&amp;file=Media+S15E12+1080p+DSNP+WEB-DL+DDP+5.1+H.264-NTb" length="558112511" type="application/x-bittorrent" />
      <torznab:attr name="category" value="5000" />
      <torznab:attr name="category" value="100002" />
      <torznab:attr name="imdb" value="1561755" />
      <torznab:attr name="imdbid" value="tt1561755" />
      <torznab:attr name="tmdbid" value="32726" />
      <torznab:attr name="seeders" value="129" />
      <torznab:attr name="peers" value="129" />
      <torznab:attr name="infohash" value="PLACEHOLDER_INFOHASH_3" />
      <torznab:attr name="downloadvolumefactor" value="1" />
      <torznab:attr name="uploadvolumefactor" value="1" />
    </item>
    <item>
      <title>Media S15E12 1080p WEB-DL DDP 5.1 H.264-ETHEL</title>
      <guid>https://generic-tracker.example/torrents/media-s15e12-1080p-web-dl-ddp-51-h264-ethel.444752</guid>
      <jackettindexer id="generic-tracker-api">Generic-Tracker (API)</jackettindexer>
      <type>private</type>
      <comments>https://generic-tracker.example/torrents/media-s15e12-1080p-web-dl-ddp-51-h264-ethel.444752</comments>
      <pubDate>Fri, 30 May 2025 07:03:54 +0000</pubDate>
      <size>541779433</size>
      <grabs>200</grabs>
      <description />
      <link>http://localhost:9117/dl/generic-tracker-api/?jackett_apikey=PLACEHOLDER_API_KEY&amp;path=PLACEHOLDER_ENCODED_PATH&amp;file=Media+S15E12+1080p+WEB-DL+DDP+5.1+H.264-ETHEL</link>
      <category>5000</category>
      <category>100002</category>
      <enclosure url="http://localhost:9117/dl/generic-tracker-api/?jackett_apikey=PLACEHOLDER_API_KEY&amp;path=PLACEHOLDER_ENCODED_PATH&amp;file=Media+S15E12+1080p+WEB-DL+DDP+5.1+H.264-ETHEL" length="541779433" type="application/x-bittorrent" />
      <torznab:attr name="category" value="5000" />
      <torznab:attr name="category" value="100002" />
      <torznab:attr name="imdb" value="1561755" />
      <torznab:attr name="imdbid" value="tt1561755" />
      <torznab:attr name="tmdbid" value="32726" />
      <torznab:attr name="seeders" value="207" />
      <torznab:attr name="peers" value="207" />
      <torznab:attr name="infohash" value="PLACEHOLDER_INFOHASH_4" />
      <torznab:attr name="downloadvolumefactor" value="1" />
      <torznab:attr name="uploadvolumefactor" value="1" />
    </item>
  </channel>
</rss>`

const exampleResponseTemplate = `
<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom" xmlns:torznab="http://torznab.com/schemas/2015/feed">
  <channel>
    <atom:link href="http://localhost:9117/" rel="self" type="application/rss+xml" />
    <title>%[1]s</title>
    <language>en-US</language>
    <category>search</category>
    <item>
      <title>A result</title>
      <guid>some-guid-string</guid>
      <jackettindexer id="%[1]s">%[1]s</jackettindexer>
      <type>private</type>
      <pubDate>Fri, 30 May 2025 22:50:38 +0000</pubDate>
      <size>39891245</size>
      <grabs>22</grabs>
      <description />
      <link>http://localhost:9117/dl/%[1]s</link>
      <category>5000</category>
      <category>100002</category>
      <enclosure url="http://localhost:9117/dl/%[1]s" length="39891245" type="application/x-bittorrent" />
      <torznab:attr name="category" value="5000" />
      <torznab:attr name="category" value="100002" />
      <torznab:attr name="imdb" value="12345" />
      <torznab:attr name="imdbid" value="tt1234" />
      <torznab:attr name="tmdbid" value="12345" />
      <torznab:attr name="seeders" value="20" />
      <torznab:attr name="peers" value="21" />
      <torznab:attr name="infohash" value="4141640e8f08d6f04df8071aa3c585e26e6328de" />
      <torznab:attr name="downloadvolumefactor" value="1" />
      <torznab:attr name="uploadvolumefactor" value="1" />
    </item>
  </channel>
</rss>`

const exampleIndexerResponse = `
<?xml version="1.0" encoding="UTF-8"?>
<indexers>
  <indexer id="my-indexer" configured="true">
    <title>My Indexer</title>
    <description>This is a fake indexer</description>
    <link>https://example.com</link>
    <language>en-US</language>
    <type>private</type>
    <caps>
      <server title="Jackett" />
      <limits default="100" max="100" />
      <searching>
        <search available="yes" supportedParams="q" />
        <tv-search available="yes" supportedParams="q,season,ep,tmdbid" />
        <movie-search available="yes" supportedParams="q,imdbid,tmdbid" />
        <music-search available="yes" supportedParams="q" />
        <audio-search available="yes" supportedParams="q" />
        <book-search available="yes" supportedParams="q" />
      </searching>
      <categories>
        <category id="2000" name="Movies" />
        <category id="5000" name="TV" />
        <category id="100001" name="Movies" />
        <category id="100002" name="TV" />
      </categories>
    </caps>
  </indexer>
</indexers>`

const exampleCapsResponse = `
<?xml version="1.0" encoding="UTF-8"?>
<caps>
  <server title="Jackett" />
  <limits default="100" max="100" />
  <searching>
    <search available="yes" supportedParams="q" />
    <tv-search available="yes" supportedParams="q,season,ep,tmdbid" />
    <movie-search available="yes" supportedParams="q,imdbid,tmdbid" />
    <music-search available="yes" supportedParams="q" />
    <audio-search available="yes" supportedParams="q" />
    <book-search available="yes" supportedParams="q" />
  </searching>
  <categories>
    <category id="2000" name="Movies" />
    <category id="5000" name="TV" />
    <category id="100001" name="Movies" />
    <category id="100002" name="TV" />
  </categories>
</caps>`

func TestCookieAuthentication(t *testing.T) {
	t.Parallel()

	// Create a mock server that checks for cookies
	var receivedCookies string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedCookies = r.Header.Get("Cookie")
		w.Write([]byte(exampleResponse))
	}))
	defer srv.Close()

	// Test cookie authentication
	j, err := New(Settings{
		ApiURL: srv.URL,
		Cookies: map[string]string{
			"Jackett": testCookieValue,
		},
		Client: srv.Client(),
	})
	if err != nil {
		t.Fatal(err)
	}

	// Make a request
	_, err = j.Fetch(context.Background(), NewRawSearch().Build())
	if err != nil {
		t.Fatal(err)
	}

	// Check that cookies were sent
	expectedCookie := fmt.Sprintf("Jackett=%s", testCookieValue)
	if receivedCookies != expectedCookie {
		t.Errorf("Expected cookie %q, got %q", expectedCookie, receivedCookies)
	}
}

func TestCookieAuthenticationWithEnvironmentVariable(t *testing.T) {
	// Set environment variable
	t.Setenv("JACKETT_COOKIES", "Jackett="+testCookieValue+"; TestCookie=1")

	// Create a mock server that checks for cookies
	var receivedCookies string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedCookies = r.Header.Get("Cookie")
		w.Write([]byte(exampleResponse))
	}))
	defer srv.Close()

	// Test cookie authentication from environment
	j, err := New(Settings{
		ApiURL: srv.URL,
		Client: srv.Client(),
	})
	if err != nil {
		t.Fatal(err)
	}

	// Make a request
	_, err = j.Fetch(context.Background(), NewRawSearch().Build())
	if err != nil {
		t.Fatal(err)
	}

	// Check that cookies were sent (order may vary)
	expectedCookies := []string{
		fmt.Sprintf("Jackett=%s; TestCookie=1", testCookieValue),
		fmt.Sprintf("TestCookie=1; Jackett=%s", testCookieValue),
	}

	found := false
	for _, expected := range expectedCookies {
		if receivedCookies == expected {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected one of %v, got %q", expectedCookies, receivedCookies)
	}
}

func TestBothAPIKeyAndCookieAuthentication(t *testing.T) {
	t.Parallel()

	// Create a mock server that checks for both API key and cookies
	var receivedAPIKey string
	var receivedCookies string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAPIKey = r.URL.Query().Get("apikey")
		receivedCookies = r.Header.Get("Cookie")
		w.Write([]byte(exampleResponse))
	}))
	defer srv.Close()

	// Test both authentication methods
	j, err := New(Settings{
		ApiURL: srv.URL,
		ApiKey: testAPIKey,
		Cookies: map[string]string{
			"Jackett": testCookieValue,
		},
		Client: srv.Client(),
	})
	if err != nil {
		t.Fatal(err)
	}

	// Make a request
	_, err = j.Fetch(context.Background(), NewRawSearch().Build())
	if err != nil {
		t.Fatal(err)
	}

	// Check that both API key and cookies were sent
	if receivedAPIKey != testAPIKey {
		t.Errorf("Expected API key %q, got %q", testAPIKey, receivedAPIKey)
	}

	expectedCookie := fmt.Sprintf("Jackett=%s", testCookieValue)
	if receivedCookies != expectedCookie {
		t.Errorf("Expected cookie %q, got %q", expectedCookie, receivedCookies)
	}
}

func TestCookieAcquisitionFlow(t *testing.T) {
	t.Parallel()

	t.Run("acquireCookies", func(t *testing.T) {
		t.Parallel()

		// Create a mock server that simulates the Jackett redirect flow
		redirectCount := 0
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/UI/Dashboard":
				if redirectCount == 0 {
					redirectCount++
					w.Header().Set("Location", "/UI/Login?ReturnUrl=%2FUI%2FDashboard")
					w.WriteHeader(http.StatusFound)
					return
				}
				// Final success response
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Dashboard"))
			case "/UI/Login":
				if r.URL.Query().Get("cookiesChecked") == "1" {
					// Final redirect back to dashboard
					w.Header().Set("Location", "/UI/Dashboard")
					w.WriteHeader(http.StatusFound)
					return
				}
				// First login redirect
				w.Header().Set("Location", "/UI/TestCookie")
				w.WriteHeader(http.StatusFound)
			case "/UI/TestCookie":
				// Set a cookie and redirect back to login
				http.SetCookie(w, &http.Cookie{
					Name:  "Jackett",
					Value: testCookieValue,
					Path:  "/",
				})
				w.Header().Set("Location", "/UI/Login?cookiesChecked=1")
				w.WriteHeader(http.StatusFound)
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		}))
		defer srv.Close()

		// Create client
		client, err := New(Settings{
			ApiURL: srv.URL,
		})
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		// Test cookie acquisition
		ctx := context.Background()
		err = client.acquireCookies(ctx)
		if err != nil {
			t.Fatalf("Failed to acquire cookies: %v", err)
		}

		// Verify cookies were stored
		client.cookieMutex.RLock()
		cookieValue, exists := client.cfg.Cookies["Jackett"]
		client.cookieMutex.RUnlock()

		if !exists {
			t.Fatal("Expected Jackett cookie to be stored")
		}
		if cookieValue != testCookieValue {
			t.Errorf("Expected cookie value %q, got %q", testCookieValue, cookieValue)
		}
	})

	t.Run("middleware retry on auth failure", func(t *testing.T) {
		t.Parallel()

		// Track request attempts
		requestCount := 0
		cookieAcquired := false

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestCount++

			// Simulate the redirect flow for cookie acquisition
			if r.URL.Path == "/UI/Dashboard" {
				if !cookieAcquired {
					// First attempt - redirect to login
					w.Header().Set("Location", "/UI/Login")
					w.WriteHeader(http.StatusFound)
					return
				}
				// After cookie acquisition - success
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Dashboard"))
				return
			}

			if r.URL.Path == "/UI/Login" {
				// Set cookie and redirect back
				http.SetCookie(w, &http.Cookie{
					Name:  "Jackett",
					Value: testCookieValue,
					Path:  "/",
				})
				w.Header().Set("Location", "/UI/Dashboard")
				w.WriteHeader(http.StatusFound)
				cookieAcquired = true
				return
			}

			// API endpoint - first call fails with 302, second succeeds
			if r.URL.Path == "/api/v2.0/indexers/all/results/torznab" {
				if requestCount <= 2 { // First API call fails
					w.Header().Set("Location", "/UI/Dashboard")
					w.WriteHeader(http.StatusFound)
					return
				}
				// Subsequent calls succeed
				w.Header().Set("Content-Type", "application/xml")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<indexers>
  <indexer id="test" title="Test Indexer" description="Test" link="http://test.com" language="en" type="public" configured="true" />
</indexers>`))
				return
			}

			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		// Create client
		client, err := New(Settings{
			ApiURL: srv.URL,
		})
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		// Make an API call that should trigger cookie acquisition and retry
		ctx := context.Background()
		_, err = client.ListIndexers(ctx, nil)
		if err != nil {
			t.Fatalf("Failed to list indexers: %v", err)
		}

		// Verify that multiple requests were made (initial + retry)
		if requestCount < 3 {
			t.Errorf("Expected at least 3 requests (cookie flow + API retry), got %d", requestCount)
		}

		// Verify cookies were stored
		client.cookieMutex.RLock()
		_, exists := client.cfg.Cookies["Jackett"]
		client.cookieMutex.RUnlock()

		if !exists {
			t.Fatal("Expected Jackett cookie to be stored after retry")
		}
	})

	t.Run("concurrent cookie acquisition", func(t *testing.T) {
		t.Parallel()

		// Track concurrent requests
		var concurrentRequests int32
		var maxConcurrent int32
		var cookieAcquired bool

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			current := atomic.AddInt32(&concurrentRequests, 1)
			defer atomic.AddInt32(&concurrentRequests, -1)

			// Track maximum concurrent requests
			for {
				currentMax := atomic.LoadInt32(&maxConcurrent)
				if current > currentMax {
					if atomic.CompareAndSwapInt32(&maxConcurrent, currentMax, current) {
						break
					}
				} else {
					break
				}
			}

			// Simulate cookie acquisition flow
			if r.URL.Path == "/UI/Dashboard" {
				http.SetCookie(w, &http.Cookie{
					Name:  "Jackett",
					Value: testCookieValue,
					Path:  "/",
				})
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Dashboard"))
				cookieAcquired = true
				return
			}

			// API endpoint - first call fails with 302 to trigger cookie acquisition
			if r.URL.Path == "/api/v2.0/indexers/all/results/torznab" {
				if !cookieAcquired {
					w.Header().Set("Location", "/UI/Dashboard")
					w.WriteHeader(http.StatusFound)
					return
				}
				w.Header().Set("Content-Type", "application/xml")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<indexers>
  <indexer id="test" title="Test Indexer" description="Test" link="http://test.com" language="en" type="public" configured="true" />
</indexers>`))
				return
			}

			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		// Create client
		client, err := New(Settings{
			ApiURL: srv.URL,
		})
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		// Make concurrent API calls
		ctx := context.Background()
		var wg sync.WaitGroup
		numGoroutines := 5

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, err := client.ListIndexers(ctx, nil)
				if err != nil {
					t.Errorf("Failed to list indexers: %v", err)
				}
			}()
		}

		wg.Wait()

		// Verify that cookie acquisition was thread-safe
		// (max concurrent should be limited by the mutex)
		if atomic.LoadInt32(&maxConcurrent) > 1 {
			t.Errorf("Expected cookie acquisition to be serialized, but saw %d concurrent requests", atomic.LoadInt32(&maxConcurrent))
		}

		// Verify cookies were stored
		client.cookieMutex.RLock()
		_, exists := client.cfg.Cookies["Jackett"]
		client.cookieMutex.RUnlock()

		if !exists {
			t.Fatal("Expected Jackett cookie to be stored")
		}
	})

	t.Run("isAuthFailure", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			statusCode int
			expected   bool
		}{
			{200, false},
			{201, false},
			{302, true},
			{401, true},
			{403, true},
			{404, false},
			{500, false},
		}

		for _, test := range tests {
			result := isAuthFailure(test.statusCode)
			if result != test.expected {
				t.Errorf("isAuthFailure(%d) = %v, expected %v", test.statusCode, result, test.expected)
			}
		}
	})
}
