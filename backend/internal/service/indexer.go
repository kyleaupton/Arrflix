package service

import (
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

// Helpers

// Get makes a GET request to the given URL and returns the response body as a []byte.
// It returns an error if the request fails or the response status is not 2xx.
func get(url string, headers map[string]string) ([]byte, error) {
	client := &http.Client{
		Timeout: 60 * time.Second, // good default timeout
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	return body, nil
}
