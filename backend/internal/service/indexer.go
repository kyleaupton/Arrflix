package service

import (
	"context"

	"github.com/kyleaupton/snaggle/backend/internal/config"
	"github.com/kyleaupton/snaggle/backend/internal/jackett"
	"github.com/kyleaupton/snaggle/backend/internal/logger"
	"github.com/kyleaupton/snaggle/backend/internal/repo"
)

type IndexerService struct {
	repo    *repo.Repository
	logger  *logger.Logger
	jackett *jackett.Client
}

func NewIndexerService(r *repo.Repository, l *logger.Logger, c *config.Config) *IndexerService {
	j, err := jackett.New(jackett.Settings{
		ApiURL: "http://localhost:9117",
		ApiKey: c.JackettAPIKey,
	})
	if err != nil {
		panic(err)
	}

	return &IndexerService{repo: r, logger: l, jackett: j}
}

// IndexersConfigured returns all configured indexers from Jackett
func (s *IndexerService) IndexersConfigured(ctx context.Context) ([]map[string]any, error) {
	s.logger.Info().Msg("Fetching configured indexers")

	indexers, err := s.jackett.ListConfiguredIndexers(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to list configured indexers")
		return nil, err
	}

	// Convert to map format
	var configured []map[string]any
	for _, indexer := range indexers {
		indexerMap := map[string]any{
			"id":          indexer.ID,
			"title":       indexer.Title,
			"description": indexer.Description,
			"link":        indexer.Link,
			"language":    indexer.Language,
			"type":        indexer.Type,
			"configured":  indexer.Configured,
		}
		configured = append(configured, indexerMap)
	}

	s.logger.Info().Int("count", len(configured)).Msg("Successfully fetched configured indexers")
	return configured, nil
}

// IndexersUnconfigured returns all unconfigured indexers from Jackett
func (s *IndexerService) IndexersUnconfigured(ctx context.Context) ([]map[string]any, error) {
	s.logger.Info().Msg("Fetching unconfigured indexers")

	indexers, err := s.jackett.ListIndexers(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to list indexers")
		return nil, err
	}

	// Filter for unconfigured indexers and convert to map format
	var unconfigured []map[string]any
	for _, indexer := range indexers {
		if indexer.Configured != "true" {
			indexerMap := map[string]any{
				"id":          indexer.ID,
				"title":       indexer.Title,
				"description": indexer.Description,
				"link":        indexer.Link,
				"language":    indexer.Language,
				"type":        indexer.Type,
				"configured":  indexer.Configured,
			}
			unconfigured = append(unconfigured, indexerMap)
		}
	}

	s.logger.Info().Int("count", len(unconfigured)).Msg("Successfully fetched unconfigured indexers")
	return unconfigured, nil
}

// AddIndexer creates a new indexer with the given configuration
func (s *IndexerService) AddIndexer(ctx context.Context, config *jackett.IndexerConfigRequest) (*jackett.IndexerConfig, error) {
	s.logger.Info().Str("name", config.Name).Msg("Creating new indexer")

	indexer, err := s.jackett.CreateIndexer(ctx, config)
	if err != nil {
		s.logger.Error().Str("name", config.Name).Err(err).Msg("Failed to create indexer")
		return nil, err
	}

	s.logger.Info().Str("id", indexer.ID).Str("name", indexer.Name).Msg("Successfully created indexer")
	return indexer, nil
}

// GetIndexerConfig retrieves the configuration for a specific indexer
func (s *IndexerService) GetIndexerConfig(ctx context.Context, indexerID string) (*jackett.IndexerConfig, error) {
	s.logger.Info().Str("indexerID", indexerID).Msg("Getting indexer configuration")

	config, err := s.jackett.GetIndexerConfig(ctx, indexerID)
	if err != nil {
		s.logger.Error().Str("indexerID", indexerID).Err(err).Msg("Failed to get indexer config")
		return nil, err
	}

	s.logger.Info().Str("indexerID", indexerID).Msg("Successfully retrieved indexer config")
	return config, nil
}

// UpdateIndexerConfig updates the configuration for a specific indexer
func (s *IndexerService) UpdateIndexerConfig(ctx context.Context, indexerID string, config *jackett.IndexerConfigRequest) (*jackett.IndexerConfig, error) {
	s.logger.Info().Str("indexerID", indexerID).Msg("Updating indexer configuration")

	updatedConfig, err := s.jackett.SaveIndexerConfig(ctx, indexerID, config)
	if err != nil {
		s.logger.Error().Str("indexerID", indexerID).Err(err).Msg("Failed to update indexer config")
		return nil, err
	}

	s.logger.Info().Str("indexerID", indexerID).Msg("Successfully updated indexer config")
	return updatedConfig, nil
}

// DeleteIndexer removes an indexer by its ID
func (s *IndexerService) DeleteIndexer(ctx context.Context, indexerID string) error {
	s.logger.Info().Str("indexerID", indexerID).Msg("Deleting indexer")

	err := s.jackett.DeleteIndexer(ctx, indexerID)
	if err != nil {
		s.logger.Error().Str("indexerID", indexerID).Err(err).Msg("Failed to delete indexer")
		return err
	}

	s.logger.Info().Str("indexerID", indexerID).Msg("Successfully deleted indexer")
	return nil
}

// ListAllIndexers returns all indexers (both configured and unconfigured)
func (s *IndexerService) ListAllIndexers(ctx context.Context) ([]map[string]any, error) {
	s.logger.Info().Msg("Fetching all indexers")

	indexers, err := s.jackett.ListIndexers(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to list indexers")
		return nil, err
	}

	// Convert to map format
	var allIndexers []map[string]any
	for _, indexer := range indexers {
		indexerMap := map[string]any{
			"id":          indexer.ID,
			"title":       indexer.Title,
			"description": indexer.Description,
			"link":        indexer.Link,
			"language":    indexer.Language,
			"type":        indexer.Type,
			"configured":  indexer.Configured,
		}
		allIndexers = append(allIndexers, indexerMap)
	}

	s.logger.Info().Int("count", len(allIndexers)).Msg("Successfully fetched all indexers")
	return allIndexers, nil
}
