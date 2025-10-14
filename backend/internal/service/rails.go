package service

import (
	"context"

	"github.com/kyleaupton/snaggle/backend/internal/model"
	"github.com/kyleaupton/snaggle/backend/internal/repo"
)

type RailsService struct {
	repo *repo.Repository
	tmdb *TmdbService
}

func NewRailsService(r *repo.Repository, tmdb *TmdbService) *RailsService {
	return &RailsService{repo: r, tmdb: tmdb}
}

func (s *RailsService) GetRails(ctx context.Context) ([]model.Rail, error) {
	rails := []model.Rail{}

	// trending movies
	moviesRes, err := s.tmdb.GetTrendingMovies(ctx)
	if err != nil {
		return nil, err
	}

	movies := []model.MovieRail{}
	for _, movie := range moviesRes.Results {
		movies = append(movies, model.MovieRail{
			TmdbID:      movie.ID,
			Title:       movie.Title,
			Overview:    movie.Overview,
			PosterPath:  movie.PosterPath,
			ReleaseDate: movie.ReleaseDate,
		})
	}

	rails = append(rails, model.Rail{
		Title:  "Trending Movies",
		Type:   "movie",
		Movies: movies,
	})

	// trending series
	seriesRes, err := s.tmdb.GetTrendingSeries(ctx)
	if err != nil {
		return nil, err
	}

	series := []model.SeriesRail{}
	for _, s := range seriesRes.Results {
		series = append(series, model.SeriesRail{
			TmdbID:      s.ID,
			Title:       s.Title,
			Overview:    s.Overview,
			PosterPath:  s.PosterPath,
			ReleaseDate: s.ReleaseDate,
		})
	}

	rails = append(rails, model.Rail{
		Title:  "Trending Series",
		Type:   "series",
		Series: series,
	})

	return rails, nil
}
