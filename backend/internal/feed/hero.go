package feed

import (
	"context"
	"sort"

	"github.com/kyleaupton/snaggle/backend/internal/model"
)

// HeroStrategy defines how to select the hero item for the feed
type HeroStrategy interface {
	SelectHero(ctx context.Context) (*model.HeroItem, string, error)
}

// BestBackdropFromTrendingStrategy selects hero from trending movies/series with best backdrop
type BestBackdropFromTrendingStrategy struct {
	tmdb TMDBClient
}

// NewBestBackdropFromTrendingStrategy creates the default hero strategy
func NewBestBackdropFromTrendingStrategy(tmdb TMDBClient) *BestBackdropFromTrendingStrategy {
	return &BestBackdropFromTrendingStrategy{tmdb: tmdb}
}

// SelectHero fetches trending content and picks the first item with a good backdrop
// Returns the HeroItem and its dedupe key (e.g., "movie:12345")
func (s *BestBackdropFromTrendingStrategy) SelectHero(ctx context.Context) (*model.HeroItem, string, error) {
	// Fetch trending movies and series
	candidates := make([]model.Title, 0, 40)

	// Get trending movies
	moviesSrc := &trendingMoviesProvider{tmdb: s.tmdb}
	movies, err := moviesSrc.Fetch(ctx, 20)
	if err == nil {
		candidates = append(candidates, movies...)
	}

	// Get trending series
	seriesSrc := &trendingSeriesProvider{tmdb: s.tmdb}
	series, err := seriesSrc.Fetch(ctx, 20)
	if err == nil {
		candidates = append(candidates, series...)
	}

	if len(candidates) == 0 {
		return nil, "", nil // No hero available
	}

	// Filter: must have backdrop and decent rating
	eligible := make([]model.Title, 0)
	for _, t := range candidates {
		if t.BackdropPath != "" && t.VoteAverage >= 6.5 {
			eligible = append(eligible, t)
		}
	}

	if len(eligible) == 0 {
		return nil, "", nil // No eligible hero
	}

	// Sort by popularity descending
	sort.Slice(eligible, func(i, j int) bool {
		return eligible[i].Popularity > eligible[j].Popularity
	})

	// Pick the most popular
	chosen := eligible[0]

	hero := &model.HeroItem{
		Title:        chosen.Title,
		Overview:     chosen.Overview,
		BackdropPath: chosen.BackdropPath,
		PosterPath:   chosen.PosterPath,
		TmdbID:       chosen.TmdbID,
		MediaType:    string(chosen.MediaType),
		// TrailerURL could be populated by fetching videos, but we'll leave it empty for now
	}

	dedupeKey := chosen.TitleKey()

	return hero, dedupeKey, nil
}
