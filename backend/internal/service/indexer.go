package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"golift.io/starr"
	"golift.io/starr/prowlarr"

	"github.com/kyleaupton/snaggle/backend/internal/config"
	"github.com/kyleaupton/snaggle/backend/internal/logger"
	"github.com/kyleaupton/snaggle/backend/internal/repo"
)

type IndexerService struct {
	repo           *repo.Repository
	logger         *logger.Logger
	prowlarrURL    string
	prowlarrAPIKey string
	prowlarr       *prowlarr.Prowlarr
}

func NewIndexerService(r *repo.Repository, l *logger.Logger, c *config.Config) *IndexerService {
	url := fmt.Sprintf("http://localhost:%s", c.ProwlarrPort)
	cfg := starr.New(c.ProwlarrAPIKey, url, 60*time.Second)
	prowlarrClient := prowlarr.New(cfg)

	return &IndexerService{repo: r, logger: l, prowlarrURL: url, prowlarrAPIKey: c.ProwlarrAPIKey, prowlarr: prowlarrClient}
}

// IndexersConfigured returns all configured indexers from Jackett
func (s *IndexerService) ListConfiguredIndexers(ctx context.Context) ([]*prowlarr.IndexerOutput, error) {
	configuredIndexers, err := s.prowlarr.GetIndexersContext(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to list configured indexers")
		return nil, err
	}

	return configuredIndexers, nil
}

// IndexersUnconfigured returns all unconfigured indexers from Jackett
func (s *IndexerService) GetSchema(ctx context.Context) ([]any, error) {
	// hand-roll schema call
	url := fmt.Sprintf("%s/api/v1/indexer/schema", s.prowlarrURL)
	schema, err := get(url, map[string]string{
		"x-api-key":    s.prowlarrAPIKey,
		"Accept":       "application/json",
		"Content-Type": "application/json",
		"User-Agent":   "Snaggle/1.0",
	})
	if err != nil {
		return nil, err
	}

	var schemaData []any
	err = json.Unmarshal(schema, &schemaData)
	if err != nil {
		return nil, err
	}

	return schemaData, nil
}

// GetIndexer returns a specific indexer by ID
func (s *IndexerService) GetIndexer(ctx context.Context, indexerID int64) (*prowlarr.IndexerOutput, error) {
	res, err := s.prowlarr.GetIndexerContext(ctx, indexerID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// UpdateIndexerConfig updates the configuration for a specific indexer
func (s *IndexerService) SaveIndexerConfig(ctx context.Context, input *prowlarr.IndexerInput) (*prowlarr.IndexerOutput, error) {
	res, err := s.prowlarr.AddIndexerContext(ctx, input)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// DeleteIndexer removes an indexer by its ID
func (s *IndexerService) DeleteIndexer(ctx context.Context, indexerID string) error {
	// s.logger.Info().Str("indexerID", indexerID).Msg("Deleting indexer")

	// err := s.jackett.DeleteIndexer(ctx, indexerID)
	// if err != nil {
	// 	s.logger.Error().Str("indexerID", indexerID).Err(err).Msg("Failed to delete indexer")
	// 	return err
	// }

	// s.logger.Info().Str("indexerID", indexerID).Msg("Successfully deleted indexer")
	return nil
}

// Search performs a search query against Prowlarr
func (s *IndexerService) Search(ctx context.Context, input prowlarr.SearchInput) ([]*prowlarr.Search, error) {
	results, err := s.prowlarr.SearchContext(ctx, input)
	if err != nil {
		s.logger.Error().Err(err).Str("query", input.Query).Msg("Failed to search Prowlarr")
		return nil, err
	}
	return results, nil
}

func (s *IndexerService) Action(ctx context.Context, actionName string, input interface{}) (any, error) {
	// hand-roll action call
	url := fmt.Sprintf("%s/api/v1/indexer/action/%s", s.prowlarrURL, actionName)

	// Convert input to JSON bytes
	inputJSON, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("marshaling input to JSON: %w", err)
	}

	// Add the missing headers that Prowlarr expects
	headers := map[string]string{
		"x-api-key":         s.prowlarrAPIKey,
		"Accept":            "application/json, text/javascript, */*; q=0.01",
		"Content-Type":      "application/json",
		"User-Agent":        "Snaggle/1.0",
		"X-Prowlarr-Client": "true",
		"X-Requested-With":  "XMLHttpRequest",
		"Origin":            s.prowlarrURL,
		"Referer":           fmt.Sprintf("%s/", s.prowlarrURL),
		"Sec-Fetch-Dest":    "empty",
		"Sec-Fetch-Mode":    "cors",
		"Sec-Fetch-Site":    "same-origin",
	}

	actionBytes, err := post(url, headers, bytes.NewReader(inputJSON))
	if err != nil {
		s.logger.Error().Err(err).Str("action", actionName).Msg("Failed to perform indexer action")
		return nil, err
	}

	// Parse the JSON response to avoid base64 encoding when serializing
	var action any
	err = json.Unmarshal(actionBytes, &action)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling action response: %w", err)
	}

	return action, nil
}

// Helpers

// httpRequest makes an HTTP request to the given URL and returns the response body as a []byte.
// It returns an error if the request fails or the response status is not 2xx.
func httpRequest(method, url string, headers map[string]string, body io.Reader) ([]byte, error) {
	client := &http.Client{
		Timeout: 60 * time.Second, // good default timeout
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	// Optional headers (can pass nil if none)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("performing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	return responseBody, nil
}

// get makes a GET request to the given URL and returns the response body as a []byte.
// It returns an error if the request fails or the response status is not 2xx.
func get(url string, headers map[string]string) ([]byte, error) {
	return httpRequest(http.MethodGet, url, headers, nil)
}

// post makes a POST request to the given URL and returns the response body as a []byte.
// It returns an error if the request fails or the response status is not 2xx.
func post(url string, headers map[string]string, body io.Reader) ([]byte, error) {
	return httpRequest(http.MethodPost, url, headers, body)
}
