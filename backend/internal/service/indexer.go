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

// ListAllIndexers returns all indexers (both configured and unconfigured)
func (s *IndexerService) ListAllIndexers(ctx context.Context) ([]jackett.IndexerDetails, error) {
	s.logger.Info().Msg("Fetching all indexers")

	indexers, err := s.jackett.ListIndexers(ctx, nil)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to list indexers")
		return nil, err
	}

	s.logger.Info().Int("count", len(indexers)).Msg("Successfully fetched all indexers")
	return indexers, nil
}

// IndexersConfigured returns all configured indexers from Jackett
func (s *IndexerService) IndexersConfigured(ctx context.Context) ([]jackett.IndexerDetails, error) {
	s.logger.Info().Msg("Fetching configured indexers")

	configured := true
	indexers, err := s.jackett.ListIndexers(ctx, &configured)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to list configured indexers")
		return nil, err
	}

	s.logger.Info().Int("count", len(indexers)).Msg("Successfully fetched configured indexers")
	return indexers, nil
}

// IndexersUnconfigured returns all unconfigured indexers from Jackett
func (s *IndexerService) IndexersUnconfigured(ctx context.Context) ([]jackett.IndexerDetails, error) {
	s.logger.Info().Msg("Fetching unconfigured indexers")

	configured := false
	indexers, err := s.jackett.ListIndexers(ctx, &configured)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to list indexers")
		return nil, err
	}

	s.logger.Info().Int("count", len(indexers)).Msg("Successfully fetched unconfigured indexers")
	return indexers, nil
}

// GetIndexerConfig retrieves the configuration for a specific indexer
func (s *IndexerService) GetIndexerConfig(ctx context.Context, indexerID string) (*jackett.IndexerConfigResponse, error) {
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
func (s *IndexerService) SaveIndexerConfig(ctx context.Context, indexerID string, config any) error {
	s.logger.Info().Str("indexerID", indexerID).Msg("Updating indexer configuration")

	err := s.jackett.SaveIndexerConfig(ctx, indexerID, config)
	if err != nil {
		s.logger.Error().Str("indexerID", indexerID).Err(err).Msg("Failed to update indexer config")
		return err
	}

	s.logger.Info().Str("indexerID", indexerID).Msg("Successfully updated indexer config")
	return nil
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
