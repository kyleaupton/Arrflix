package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	tmdb "github.com/cyruzin/golang-tmdb"
	"github.com/kyleaupton/snaggle/backend/internal/config"
	"github.com/kyleaupton/snaggle/backend/internal/logger"
	"github.com/kyleaupton/snaggle/backend/internal/repo"
)

// TTL for static data (7 days)
const STATIC_TTL = (24 * time.Hour) * 7

// TTL for dynamic data (1 hour)
const DYNAMIC_TTL = time.Hour

type TmdbService struct {
	repo   *repo.Repository
	client *tmdb.Client
	logger *logger.Logger
}

func NewTmdbService(r *repo.Repository, l *logger.Logger) *TmdbService {
	client, err := tmdb.Init(config.Load().TmdbAPIKey)
	if err != nil {
		panic(err)
	}

	return &TmdbService{
		repo:   r,
		logger: l,
		client: client,
	}
}

func (s *TmdbService) FindByID(ctx context.Context, id, source string) (tmdb.FindByID, error) {
	cacheKey := fmt.Sprintf("tmdb_find_by_id_%s_%s", id, source)
	return getOrFetchFromCache(ctx, s.repo, s.logger, cacheKey, func() (*tmdb.FindByID, error) {
		return s.client.GetFindByID(id, map[string]string{
			"external_source": source,
		})
	}, STATIC_TTL)
}

func (s *TmdbService) GetMovieDetails(ctx context.Context, id int64) (tmdb.MovieDetails, error) {
	cacheKey := fmt.Sprintf("tmdb_movie_details_%d", id)
	return getOrFetchFromCache(ctx, s.repo, s.logger, cacheKey, func() (*tmdb.MovieDetails, error) {
		return s.client.GetMovieDetails(int(id), map[string]string{})
	}, STATIC_TTL)
}

func (s *TmdbService) GetSeriesDetails(ctx context.Context, id int64) (tmdb.TVDetails, error) {
	cacheKey := fmt.Sprintf("tmdb_series_details_%d", id)
	return getOrFetchFromCache(ctx, s.repo, s.logger, cacheKey, func() (*tmdb.TVDetails, error) {
		return s.client.GetTVDetails(int(id), map[string]string{})
	}, STATIC_TTL)
}

func (s *TmdbService) GetEpisodeDetails(ctx context.Context, id int64, season int64, episode int64) (tmdb.TVEpisodeDetails, error) {
	cacheKey := fmt.Sprintf("tmdb_episode_details_%d_%d_%d", id, season, episode)
	return getOrFetchFromCache(ctx, s.repo, s.logger, cacheKey, func() (*tmdb.TVEpisodeDetails, error) {
		return s.client.GetTVEpisodeDetails(int(id), int(season), int(episode), map[string]string{})
	}, STATIC_TTL)
}

func (s *TmdbService) GetTrendingMovies(ctx context.Context) (tmdb.Trending, error) {
	cacheKey := "tmdb_trending_movies"
	return getOrFetchFromCache(ctx, s.repo, s.logger, cacheKey, func() (*tmdb.Trending, error) {
		return s.client.GetTrending("movie", "day", map[string]string{})
	}, DYNAMIC_TTL)
}

func (s *TmdbService) GetTrendingSeries(ctx context.Context) (tmdb.Trending, error) {
	cacheKey := "tmdb_trending_series"
	return getOrFetchFromCache(ctx, s.repo, s.logger, cacheKey, func() (*tmdb.Trending, error) {
		return s.client.GetTrending("tv", "day", map[string]string{})
	}, DYNAMIC_TTL)
}

// getOrFetchFromCache encapsulates the pattern of:
// 1) checking API cache
// 2) calling the provided fetch function on cache miss
// 3) storing the fresh response back into the cache
// 4) returning the typed result
func getOrFetchFromCache[T any](ctx context.Context, r *repo.Repository, l *logger.Logger, cacheKey string, fetch func() (*T, error), ttl time.Duration) (T, error) {
	cacheEntry, found, err := r.GetApiCache(ctx, cacheKey)
	if err != nil {
		var zero T
		return zero, err
	}

	if !found {
		l.Debug().Str("cache_key", cacheKey).Msg("Cache miss, fetching from API")
		res, err := fetch()
		if err != nil {
			var zero T
			return zero, err
		}

		category := "tmdb"
		contentType := "application/json"

		// Convert the result to json to be stored in the cache
		jsonRes, err := json.Marshal(res)
		if err != nil {
			var zero T
			return zero, err
		}

		// Note: pass nil for headers so the DB receives NULL (valid for jsonb)
		if err := r.UpsertApiCache(ctx, cacheKey, &category, jsonRes, 200, &contentType, nil, ttl); err != nil {
			l.Error().Err(err).Str("cache_key", cacheKey).Msg("Failed upserting api cache")
		}

		// Return the result
		return *res, nil
	}

	var out T
	err = json.Unmarshal(cacheEntry.Response, &out)
	if err != nil {
		var zero T
		return zero, err
	}
	return out, nil
}
