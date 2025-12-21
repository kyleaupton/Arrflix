package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"golift.io/starr/prowlarr"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	dbgen "github.com/kyleaupton/snaggle/backend/internal/db/sqlc"
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
		candidate := s.searchResultToCandidate(result)
		candidates = append(candidates, candidate)
	}

	return candidates, nil
}

// EvaluateCandidate returns the evaluation trace for a candidate
func (s *DownloadCandidatesService) EvaluateCandidate(ctx context.Context, movieID int64, indexerID int64, guid string) (model.EvaluationTrace, error) {
	// Lookup torrent from cache
	cacheKey := s.cacheKey(indexerID, guid)
	s.cacheMu.RLock()
	cached, exists := s.cache[cacheKey]
	s.cacheMu.RUnlock()

	if !exists {
		return model.EvaluationTrace{}, ErrCandidateNotFound
	}

	// Check if expired
	if time.Now().After(cached.expiresAt) {
		s.cacheMu.Lock()
		delete(s.cache, cacheKey)
		s.cacheMu.Unlock()
		return model.EvaluationTrace{}, ErrCandidateExpired
	}

	result := cached.result

	// Transform to DownloadCandidate
	candidate := s.searchResultToCandidate(result)

	trace, err := s.policyEngine.Evaluate(ctx, candidate)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to evaluate policy")
		return model.EvaluationTrace{}, fmt.Errorf("failed to evaluate policy: %w", err)
	}

	return trace, nil
}

// PreviewCandidate previews what will happen when a candidate is enqueued.
func (s *DownloadCandidatesService) PreviewCandidate(ctx context.Context, movieID int64, indexerID int64, guid string) (model.EvaluationTrace, error) {
	return s.EvaluateCandidate(ctx, movieID, indexerID, guid)
}

// EnqueueCandidate creates a durable download job for a candidate (movies-only for v1).
func (s *DownloadCandidatesService) EnqueueCandidate(ctx context.Context, movieID int64, indexerID int64, guid string) (model.EvaluationTrace, dbgen.DownloadJob, error) {
	// Lookup candidate from cache
	cacheKey := s.cacheKey(indexerID, guid)
	s.cacheMu.RLock()
	cached, exists := s.cache[cacheKey]
	s.cacheMu.RUnlock()

	if !exists {
		return model.EvaluationTrace{}, dbgen.DownloadJob{}, ErrCandidateNotFound
	}

	// Check if expired
	if time.Now().After(cached.expiresAt) {
		s.cacheMu.Lock()
		delete(s.cache, cacheKey)
		s.cacheMu.Unlock()
		return model.EvaluationTrace{}, dbgen.DownloadJob{}, ErrCandidateExpired
	}

	candidate := s.searchResultToCandidate(cached.result)

	s.logger.Info().Interface("candidate", candidate).Msg("Candidate")

	trace, err := s.policyEngine.Evaluate(ctx, candidate)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to evaluate policy")
		return model.EvaluationTrace{}, dbgen.DownloadJob{}, fmt.Errorf("failed to evaluate policy: %w", err)
	}

	var downloaderID, libraryID, nameTemplateID pgtype.UUID
	if err := downloaderID.Scan(trace.FinalPlan.DownloaderID); err != nil {
		return trace, dbgen.DownloadJob{}, fmt.Errorf("invalid downloader id: %w", err)
	}
	if err := libraryID.Scan(trace.FinalPlan.LibraryID); err != nil {
		return trace, dbgen.DownloadJob{}, fmt.Errorf("invalid library id: %w", err)
	}
	if err := nameTemplateID.Scan(trace.FinalPlan.NameTemplateID); err != nil {
		return trace, dbgen.DownloadJob{}, fmt.Errorf("invalid name template id: %w", err)
	}

	// Ensure media_item exists for this movie/library and link the job to it.
	mi, err := s.repo.GetMediaItemByTmdbID(ctx, movieID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			movie, err := s.media.GetMovie(ctx, movieID)
			if err != nil {
				return trace, dbgen.DownloadJob{}, fmt.Errorf("get movie: %w", err)
			}
			var yearInt *int32
			if len(movie.ReleaseDate) >= 4 {
				if y, err := strconv.Atoi(movie.ReleaseDate[:4]); err == nil {
					yy := int32(y)
					yearInt = &yy
				}
			}
			tmdb := movieID
			mi, err = s.repo.CreateMediaItem(ctx, libraryID, "movie", movie.Title, yearInt, &tmdb)
			if err != nil {
				return trace, dbgen.DownloadJob{}, fmt.Errorf("create media item: %w", err)
			}
		} else {
			return trace, dbgen.DownloadJob{}, fmt.Errorf("get media item: %w", err)
		}
	}

	job, err := s.repo.CreateDownloadJob(ctx, dbgen.CreateDownloadJobParams{
		Protocol:       candidate.Protocol,
		MediaType:      "movie",
		MediaItemID:    mi.ID,
		SeasonID:       pgtype.UUID{},
		EpisodeID:      pgtype.UUID{},
		IndexerID:      indexerID,
		Guid:           guid,
		CandidateTitle: candidate.Title,
		CandidateLink:  candidate.Link,
		DownloaderID:   downloaderID,
		LibraryID:      libraryID,
		NameTemplateID: nameTemplateID,
	})
	if err != nil {
		return trace, dbgen.DownloadJob{}, fmt.Errorf("create download job: %w", err)
	}

	return trace, job, nil
}

// searchResultToCandidate converts a prowlarr.Search result to a model.DownloadCandidate
func (s *DownloadCandidatesService) searchResultToCandidate(result *prowlarr.Search) model.DownloadCandidate {
	// Extract categories
	categories := make([]string, 0, len(result.Categories))
	for _, cat := range result.Categories {
		if cat != nil && cat.Name != "" {
			categories = append(categories, cat.Name)
		}
	}

	return model.DownloadCandidate{
		Protocol: string(result.Protocol),
		Filename: result.FileName,
		// Link:        result.DownloadURL,
		Link:        result.GUID,
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
