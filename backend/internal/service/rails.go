package service

import (
	"context"
	"strconv"

	tmdb "github.com/cyruzin/golang-tmdb"
	dbgen "github.com/kyleaupton/snaggle/backend/internal/db/sqlc"
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

	svc := s

	// Helper function to convert trending movies to MovieRail
	convertTrendingMovies := func(moviesRes tmdb.Trending) []model.MovieRail {
		movies := []model.MovieRail{}
		for _, item := range moviesRes.Results {
			year := parseYearFromDate(item.ReleaseDate)
			movies = append(movies, model.MovieRail{
				TmdbID:        item.ID,
				Title:         item.Title,
				Overview:      item.Overview,
				PosterPath:    item.PosterPath,
				ReleaseDate:   item.ReleaseDate,
				Year:          year,
				Genres:        toInt64FromInt(item.GenreIDs),
				Tagline:       "",
				IsInLibrary:   svc.isInLibrary(ctx, item.ID, model.MediaTypeMovie),
				IsDownloading: svc.hasActiveDownloads(ctx, item.ID, model.MediaTypeMovie),
			})
		}
		return movies
	}

	// Helper function to convert trending series to SeriesRail
	convertTrendingSeries := func(seriesRes tmdb.Trending) []model.SeriesRail {
		series := []model.SeriesRail{}
		for _, item := range seriesRes.Results {
			year := parseYearFromDate(item.ReleaseDate)
			series = append(series, model.SeriesRail{
				TmdbID:        item.ID,
				Title:         item.Title,
				Overview:      item.Overview,
				PosterPath:    item.PosterPath,
				ReleaseDate:   item.ReleaseDate,
				Year:          year,
				Genres:        toInt64FromInt(item.GenreIDs),
				Tagline:       "",
				IsInLibrary:   svc.isInLibrary(ctx, item.ID, model.MediaTypeSeries),
				IsDownloading: svc.hasActiveDownloads(ctx, item.ID, model.MediaTypeSeries),
			})
		}
		return series
	}

	// Helper function to convert popular movies to MovieRail
	convertPopularMovies := func(moviesRes tmdb.MoviePopular) []model.MovieRail {
		movies := []model.MovieRail{}
		for _, movie := range moviesRes.Results {
			year := parseYearFromDate(movie.ReleaseDate)
			movies = append(movies, model.MovieRail{
				TmdbID:        movie.ID,
				Title:         movie.Title,
				Overview:      movie.Overview,
				PosterPath:    movie.PosterPath,
				ReleaseDate:   movie.ReleaseDate,
				Year:          year,
				Genres:        genreIDsFromGenres(movie.Genres),
				Tagline:       "",
				IsInLibrary:   svc.isInLibrary(ctx, movie.ID, model.MediaTypeMovie),
				IsDownloading: svc.hasActiveDownloads(ctx, movie.ID, model.MediaTypeMovie),
			})
		}
		return movies
	}

	// Helper function to convert top rated movies to MovieRail
	convertTopRatedMovies := func(moviesRes tmdb.MovieTopRated) []model.MovieRail {
		movies := []model.MovieRail{}
		for _, movie := range moviesRes.Results {
			year := parseYearFromDate(movie.ReleaseDate)
			movies = append(movies, model.MovieRail{
				TmdbID:        movie.ID,
				Title:         movie.Title,
				Overview:      movie.Overview,
				PosterPath:    movie.PosterPath,
				ReleaseDate:   movie.ReleaseDate,
				Year:          year,
				Genres:        genreIDsFromGenres(movie.Genres),
				Tagline:       "",
				IsInLibrary:   svc.isInLibrary(ctx, movie.ID, model.MediaTypeMovie),
				IsDownloading: svc.hasActiveDownloads(ctx, movie.ID, model.MediaTypeMovie),
			})
		}
		return movies
	}

	// Helper function to convert now playing movies to MovieRail
	convertNowPlayingMovies := func(moviesRes tmdb.MovieNowPlaying) []model.MovieRail {
		movies := []model.MovieRail{}
		for _, movie := range moviesRes.Results {
			year := parseYearFromDate(movie.ReleaseDate)
			movies = append(movies, model.MovieRail{
				TmdbID:        movie.ID,
				Title:         movie.Title,
				Overview:      movie.Overview,
				PosterPath:    movie.PosterPath,
				ReleaseDate:   movie.ReleaseDate,
				Year:          year,
				Genres:        genreIDsFromGenres(movie.Genres),
				Tagline:       "",
				IsInLibrary:   svc.isInLibrary(ctx, movie.ID, model.MediaTypeMovie),
				IsDownloading: svc.hasActiveDownloads(ctx, movie.ID, model.MediaTypeMovie),
			})
		}
		return movies
	}

	// Helper function to convert upcoming movies to MovieRail
	convertUpcomingMovies := func(moviesRes tmdb.MovieUpcoming) []model.MovieRail {
		movies := []model.MovieRail{}
		for _, movie := range moviesRes.Results {
			year := parseYearFromDate(movie.ReleaseDate)
			movies = append(movies, model.MovieRail{
				TmdbID:        movie.ID,
				Title:         movie.Title,
				Overview:      movie.Overview,
				PosterPath:    movie.PosterPath,
				ReleaseDate:   movie.ReleaseDate,
				Year:          year,
				Genres:        genreIDsFromGenres(movie.Genres),
				Tagline:       "",
				IsInLibrary:   svc.isInLibrary(ctx, movie.ID, model.MediaTypeMovie),
				IsDownloading: svc.hasActiveDownloads(ctx, movie.ID, model.MediaTypeMovie),
			})
		}
		return movies
	}

	// Helper function to convert popular series to SeriesRail
	convertPopularSeries := func(seriesRes tmdb.TVPopular) []model.SeriesRail {
		series := []model.SeriesRail{}
		for _, item := range seriesRes.Results {
			year := parseYearFromDate(item.FirstAirDate)
			series = append(series, model.SeriesRail{
				TmdbID:        item.ID,
				Title:         item.Name,
				Overview:      item.Overview,
				PosterPath:    item.PosterPath,
				ReleaseDate:   item.FirstAirDate,
				Year:          year,
				Genres:        toInt64FromInt(item.GenreIDs),
				Tagline:       "",
				IsInLibrary:   svc.isInLibrary(ctx, item.ID, model.MediaTypeSeries),
				IsDownloading: svc.hasActiveDownloads(ctx, item.ID, model.MediaTypeSeries),
			})
		}
		return series
	}

	// Helper function to convert top rated series to SeriesRail
	convertTopRatedSeries := func(seriesRes tmdb.TVTopRated) []model.SeriesRail {
		series := []model.SeriesRail{}
		for _, item := range seriesRes.Results {
			year := parseYearFromDate(item.FirstAirDate)
			series = append(series, model.SeriesRail{
				TmdbID:        item.ID,
				Title:         item.Name,
				Overview:      item.Overview,
				PosterPath:    item.PosterPath,
				ReleaseDate:   item.FirstAirDate,
				Year:          year,
				Genres:        toInt64FromInt(item.GenreIDs),
				Tagline:       "",
				IsInLibrary:   svc.isInLibrary(ctx, item.ID, model.MediaTypeSeries),
				IsDownloading: svc.hasActiveDownloads(ctx, item.ID, model.MediaTypeSeries),
			})
		}
		return series
	}

	// Helper function to convert on the air series to SeriesRail
	convertOnTheAirSeries := func(seriesRes tmdb.TVOnTheAir) []model.SeriesRail {
		series := []model.SeriesRail{}
		for _, item := range seriesRes.Results {
			year := parseYearFromDate(item.FirstAirDate)
			series = append(series, model.SeriesRail{
				TmdbID:        item.ID,
				Title:         item.Name,
				Overview:      item.Overview,
				PosterPath:    item.PosterPath,
				ReleaseDate:   item.FirstAirDate,
				Year:          year,
				Genres:        toInt64FromInt(item.GenreIDs),
				Tagline:       "",
				IsInLibrary:   svc.isInLibrary(ctx, item.ID, model.MediaTypeSeries),
				IsDownloading: svc.hasActiveDownloads(ctx, item.ID, model.MediaTypeSeries),
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

func (s *RailsService) isInLibrary(ctx context.Context, tmdbID int64, typ model.MediaType) bool {
	_, err := s.repo.GetMediaItemByTmdbIDAndType(ctx, tmdbID, string(typ))
	return err == nil
}

func (s *RailsService) hasActiveDownloads(ctx context.Context, tmdbID int64, mediaType model.MediaType) bool {
	var jobs []dbgen.DownloadJob
	var err error

	if mediaType == model.MediaTypeMovie {
		jobs, err = s.repo.ListDownloadJobsByTmdbMovieID(ctx, tmdbID)
	} else if mediaType == model.MediaTypeSeries {
		jobs, err = s.repo.ListDownloadJobsByTmdbSeriesID(ctx, tmdbID)
	} else {
		return false
	}

	if err != nil {
		return false
	}

	activeStatuses := map[string]bool{
		"created":    true,
		"enqueued":   true,
		"downloading": true,
		"importing":  true,
	}

	for _, job := range jobs {
		if activeStatuses[job.Status] {
			return true
		}
	}

	return false
}

func toInt64Slice(ints []int64) []int64 {
	if len(ints) == 0 {
		return nil
	}
	out := make([]int64, len(ints))
	copy(out, ints)
	return out
}

func toInt64FromInt(ints []int64) []int64 {
	if len(ints) == 0 {
		return nil
	}
	out := make([]int64, len(ints))
	copy(out, ints)
	return out
}

func toInt64FromInt32(ints []int) []int64 {
	if len(ints) == 0 {
		return nil
	}
	out := make([]int64, len(ints))
	for i, v := range ints {
		out[i] = int64(v)
	}
	return out
}

func genreIDsFromGenres(genres []tmdb.Genre) []int64 {
	if len(genres) == 0 {
		return nil
	}
	out := make([]int64, len(genres))
	for i, g := range genres {
		out[i] = int64(g.ID)
	}
	return out
}

func parseYearFromDate(dateStr string) *int32 {
	if len(dateStr) < 4 {
		return nil
	}
	y, err := strconv.Atoi(dateStr[:4])
	if err != nil {
		return nil
	}
	yy := int32(y)
	return &yy
}
