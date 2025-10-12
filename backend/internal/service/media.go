package service

import (
	"context"

	dbgen "github.com/kyleaupton/snaggle/backend/internal/db/sqlc"
	"github.com/kyleaupton/snaggle/backend/internal/logger"
	"github.com/kyleaupton/snaggle/backend/internal/model"
	"github.com/kyleaupton/snaggle/backend/internal/repo"
)

type MediaService struct {
	repo   *repo.Repository
	logger *logger.Logger
	tmdb   *TmdbService
}

func NewMediaService(r *repo.Repository, l *logger.Logger, tmdb *TmdbService) *MediaService {
	return &MediaService{repo: r, logger: l, tmdb: tmdb}
}

func (s *MediaService) ListLibraryItems(ctx context.Context) ([]dbgen.MediaItem, error) {
	return s.repo.ListMediaItems(ctx)
}

func (s *MediaService) GetMovie(ctx context.Context, id int64) (model.Movie, error) {
	tmdbDetails, err := s.tmdb.GetMovieDetails(ctx, id)
	if err != nil {
		return model.Movie{}, err
	}

	return model.Movie{
		TmdbID:      tmdbDetails.ID,
		Title:       tmdbDetails.Title,
		ReleaseDate: tmdbDetails.ReleaseDate,
	}, nil
}

func (s *MediaService) GetSeries(ctx context.Context, id int64) (model.Series, error) {
	tmdbDetails, err := s.tmdb.GetSeriesDetails(ctx, id)
	if err != nil {
		return model.Series{}, err
	}

	return model.Series{
		TmdbID:      tmdbDetails.ID,
		Title:       tmdbDetails.Name,
		ReleaseDate: tmdbDetails.FirstAirDate,
	}, nil
}
