package service

import (
	"context"

	tmdb "github.com/cyruzin/golang-tmdb"
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

	// Helper function to convert trending movies to MovieRail
	convertTrendingMovies := func(moviesRes tmdb.Trending) []model.MovieRail {
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
		return movies
	}

	// Helper function to convert trending series to SeriesRail
	convertTrendingSeries := func(seriesRes tmdb.Trending) []model.SeriesRail {
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
		return series
	}

	// Helper function to convert popular movies to MovieRail
	convertPopularMovies := func(moviesRes tmdb.MoviePopular) []model.MovieRail {
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
		return movies
	}

	// Helper function to convert top rated movies to MovieRail
	convertTopRatedMovies := func(moviesRes tmdb.MovieTopRated) []model.MovieRail {
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
		return movies
	}

	// Helper function to convert now playing movies to MovieRail
	convertNowPlayingMovies := func(moviesRes tmdb.MovieNowPlaying) []model.MovieRail {
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
		return movies
	}

	// Helper function to convert upcoming movies to MovieRail
	convertUpcomingMovies := func(moviesRes tmdb.MovieUpcoming) []model.MovieRail {
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
		return movies
	}

	// Helper function to convert popular series to SeriesRail
	convertPopularSeries := func(seriesRes tmdb.TVPopular) []model.SeriesRail {
		series := []model.SeriesRail{}
		for _, s := range seriesRes.Results {
			series = append(series, model.SeriesRail{
				TmdbID:      s.ID,
				Title:       s.Name,
				Overview:    s.Overview,
				PosterPath:  s.PosterPath,
				ReleaseDate: s.FirstAirDate,
			})
		}
		return series
	}

	// Helper function to convert top rated series to SeriesRail
	convertTopRatedSeries := func(seriesRes tmdb.TVTopRated) []model.SeriesRail {
		series := []model.SeriesRail{}
		for _, s := range seriesRes.Results {
			series = append(series, model.SeriesRail{
				TmdbID:      s.ID,
				Title:       s.Name,
				Overview:    s.Overview,
				PosterPath:  s.PosterPath,
				ReleaseDate: s.FirstAirDate,
			})
		}
		return series
	}

	// Helper function to convert on the air series to SeriesRail
	convertOnTheAirSeries := func(seriesRes tmdb.TVOnTheAir) []model.SeriesRail {
		series := []model.SeriesRail{}
		for _, s := range seriesRes.Results {
			series = append(series, model.SeriesRail{
				TmdbID:      s.ID,
				Title:       s.Name,
				Overview:    s.Overview,
				PosterPath:  s.PosterPath,
				ReleaseDate: s.FirstAirDate,
			})
		}
		return series
	}

	// Trending Movies
	if moviesRes, err := s.tmdb.GetTrendingMovies(ctx); err == nil {
		rails = append(rails, model.Rail{
			Title:  "Trending Movies",
			Type:   "movie",
			Movies: convertTrendingMovies(moviesRes),
		})
	}

	// Popular Movies
	if moviesRes, err := s.tmdb.GetPopularMovies(ctx); err == nil {
		rails = append(rails, model.Rail{
			Title:  "Popular Movies",
			Type:   "movie",
			Movies: convertPopularMovies(moviesRes),
		})
	}

	// Top Rated Movies
	if moviesRes, err := s.tmdb.GetTopRatedMovies(ctx); err == nil {
		rails = append(rails, model.Rail{
			Title:  "Top Rated Movies",
			Type:   "movie",
			Movies: convertTopRatedMovies(moviesRes),
		})
	}

	// Now Playing Movies
	if moviesRes, err := s.tmdb.GetNowPlayingMovies(ctx); err == nil {
		rails = append(rails, model.Rail{
			Title:  "Now Playing",
			Type:   "movie",
			Movies: convertNowPlayingMovies(moviesRes),
		})
	}

	// Upcoming Movies
	if moviesRes, err := s.tmdb.GetUpcomingMovies(ctx); err == nil {
		rails = append(rails, model.Rail{
			Title:  "Upcoming Movies",
			Type:   "movie",
			Movies: convertUpcomingMovies(moviesRes),
		})
	}

	// Trending Series
	if seriesRes, err := s.tmdb.GetTrendingSeries(ctx); err == nil {
		rails = append(rails, model.Rail{
			Title:  "Trending Series",
			Type:   "series",
			Series: convertTrendingSeries(seriesRes),
		})
	}

	// Popular Series
	if seriesRes, err := s.tmdb.GetPopularSeries(ctx); err == nil {
		rails = append(rails, model.Rail{
			Title:  "Popular Series",
			Type:   "series",
			Series: convertPopularSeries(seriesRes),
		})
	}

	// Top Rated Series
	if seriesRes, err := s.tmdb.GetTopRatedSeries(ctx); err == nil {
		rails = append(rails, model.Rail{
			Title:  "Top Rated Series",
			Type:   "series",
			Series: convertTopRatedSeries(seriesRes),
		})
	}

	// On The Air Series
	if seriesRes, err := s.tmdb.GetOnTheAirSeries(ctx); err == nil {
		rails = append(rails, model.Rail{
			Title:  "On The Air",
			Type:   "series",
			Series: convertOnTheAirSeries(seriesRes),
		})
	}

	return rails, nil
}
