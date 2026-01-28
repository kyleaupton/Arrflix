package service

import (
	"context"

	"github.com/kyleaupton/arrflix/internal/feed"
	"github.com/kyleaupton/arrflix/internal/logger"
	"github.com/kyleaupton/arrflix/internal/model"
	"github.com/kyleaupton/arrflix/internal/repo"
)

type FeedService struct {
	composer  *feed.Composer
	freshness *feed.InMemoryFreshnessTracker
	logger    *logger.Logger
}

func NewFeedService(r *repo.Repository, l *logger.Logger, tmdb *TmdbService) *FeedService {
	// Initialize feed components
	registry := feed.NewRegistry()
	sources := feed.NewTMDBSourceFactory(tmdb)
	heroStrategy := feed.NewBestBackdropFromTrendingStrategy(tmdb)
	freshnessTracker := feed.NewInMemoryFreshnessTracker()

	composer := feed.NewComposer(registry, sources, heroStrategy, r, freshnessTracker)

	return &FeedService{
		composer:  composer,
		freshness: freshnessTracker,
		logger:    l,
	}
}

func (s *FeedService) GetFeed(ctx context.Context) (*model.Feed, error) {
	return s.composer.BuildFeed(ctx)
}
