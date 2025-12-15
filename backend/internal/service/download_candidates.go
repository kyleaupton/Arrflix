package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"golift.io/starr/prowlarr"

	"github.com/kyleaupton/snaggle/backend/internal/logger"
	"github.com/kyleaupton/snaggle/backend/internal/model"
	"github.com/kyleaupton/snaggle/backend/internal/policy"
	"github.com/kyleaupton/snaggle/backend/internal/repo"
)

var (
	ErrCandidateNotFound = errors.New("candidate not found in cache (may have expired)")
	ErrCandidateExpired  = errors.New("candidate cache expired")
)

const cacheTTL = 5 * time.Minute

// cachedSearchResult stores a search result with its expiration time
type cachedSearchResult struct {
	result    *prowlarr.Search
	expiresAt time.Time
}

// DownloadCandidatesService handles download candidate search and enqueueing
type DownloadCandidatesService struct {
	repo         *repo.Repository
	logger       *logger.Logger
	indexer      *IndexerService
	media        *MediaService
	policyEngine *policy.Engine
	cache        map[string]*cachedSearchResult
	cacheMu      sync.RWMutex
}

// NewDownloadCandidatesService creates a new download candidates service
func NewDownloadCandidatesService(r *repo.Repository, l *logger.Logger, indexer *IndexerService, media *MediaService, engine *policy.Engine) *DownloadCandidatesService {
	return &DownloadCandidatesService{
		repo:         r,
		logger:       l,
		indexer:      indexer,
		media:        media,
		policyEngine: engine,
		cache:        make(map[string]*cachedSearchResult),
	}
}

// SearchDownloadCandidates searches for download candidates for a movie
func (s *DownloadCandidatesService) SearchDownloadCandidates(ctx context.Context, movieID int64) ([]model.DownloadCandidate, error) {
	// Get movie details to construct search query
	movie, err := s.media.GetMovie(ctx, movieID)
	if err != nil {
		return nil, fmt.Errorf("failed to get movie: %w", err)
	}

	// Construct search query: "Title Year"
	year := ""
	if movie.ReleaseDate != "" {
		// Extract year from release date (format: "2006-01-02")
		if len(movie.ReleaseDate) >= 4 {
			year = movie.ReleaseDate[:4]
		}
	}
	query := movie.Title
	if year != "" {
		query = fmt.Sprintf("%s %s", movie.Title, year)
	}

	// Search Prowlarr
	searchInput := prowlarr.SearchInput{
		Query: query,
		Type:  "search",
		Limit: 100,
	}

	results, err := s.indexer.Search(ctx, searchInput)
	if err != nil {
		s.logger.Error().Err(err).Str("query", query).Msg("Failed to search Prowlarr")
		return nil, fmt.Errorf("failed to search Prowlarr: %w", err)
	}

	// Clear expired entries from cache
	s.cleanExpiredCache()

	// Store results in cache and transform to DownloadCandidate
	candidates := make([]model.DownloadCandidate, 0, len(results))
	for _, result := range results {
		// Cache the result
		cacheKey := s.cacheKey(result.IndexerID, result.GUID)
		s.cacheMu.Lock()
		s.cache[cacheKey] = &cachedSearchResult{
			result:    result,
			expiresAt: time.Now().Add(cacheTTL),
		}
		s.cacheMu.Unlock()

		// Transform to DownloadCandidate
		categories := make([]string, 0, len(result.Categories))
		for _, cat := range result.Categories {
			if cat != nil && cat.Name != "" {
				categories = append(categories, cat.Name)
			}
		}

		candidate := model.DownloadCandidate{
			Protocol:    string(result.Protocol),
			Filename:    result.FileName,
			Link:        result.DownloadURL,
			Indexer:     result.Indexer,
			IndexerID:   result.IndexerID,
			GUID:        result.GUID,
			Peers:       result.Leechers,
			Seeders:     result.Seeders,
			Age:         result.Age,
			AgeHours:    result.AgeHours,
			Size:        result.Size,
			Grabs:       result.Grabs,
			Categories:  categories,
			PublishDate: result.PublishDate,
			Title:       result.Title,
		}
		candidates = append(candidates, candidate)
	}

	return candidates, nil
}

// EnqueueCandidate enqueues a download candidate through the policy engine
func (s *DownloadCandidatesService) EnqueueCandidate(ctx context.Context, movieID int64, indexerID int64, guid string) (model.Plan, error) {
	// Lookup torrent from cache
	cacheKey := s.cacheKey(indexerID, guid)
	s.cacheMu.RLock()
	cached, exists := s.cache[cacheKey]
	s.cacheMu.RUnlock()

	if !exists {
		return model.Plan{}, ErrCandidateNotFound
	}

	// Check if expired
	if time.Now().After(cached.expiresAt) {
		s.cacheMu.Lock()
		delete(s.cache, cacheKey)
		s.cacheMu.Unlock()
		return model.Plan{}, ErrCandidateExpired
	}

	result := cached.result

	// Extract categories
	categories := make([]string, 0, len(result.Categories))
	for _, cat := range result.Categories {
		if cat != nil && cat.Name != "" {
			categories = append(categories, cat.Name)
		}
	}

	// Build torrent metadata
	metadata := model.TorrentMetadata{
		Size:       uint64(result.Size),
		Seeders:    uint(result.Seeders),
		Peers:      uint(result.Leechers),
		Title:      result.Title,
		Tracker:    result.Indexer,
		TrackerID:  fmt.Sprintf("%d", result.IndexerID),
		Categories: categories,
	}

	// Evaluate through policy engine
	plan, err := s.policyEngine.Evaluate(ctx, model.EvaluateParams{
		TorrentURL: result.DownloadURL,
		Metadata:   metadata,
		MediaType:  model.MediaTypeMovie,
	})
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to evaluate policy")
		return model.Plan{}, fmt.Errorf("failed to evaluate policy: %w", err)
	}

	return plan, nil
}

// cacheKey generates a cache key from indexer ID and GUID
func (s *DownloadCandidatesService) cacheKey(indexerID int64, guid string) string {
	return fmt.Sprintf("%d:%s", indexerID, guid)
}

// cleanExpiredCache removes expired entries from the cache
func (s *DownloadCandidatesService) cleanExpiredCache() {
	now := time.Now()
	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()

	for key, cached := range s.cache {
		if now.After(cached.expiresAt) {
			delete(s.cache, key)
		}
	}
}
