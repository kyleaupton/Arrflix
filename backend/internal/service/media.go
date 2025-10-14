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

	productionCompanies := []model.ProductionCompany{}
	for _, company := range tmdbDetails.ProductionCompanies {
		productionCompanies = append(productionCompanies, model.ProductionCompany{
			TmdbID:        company.ID,
			Name:          company.Name,
			LogoPath:      company.LogoPath,
			OriginCountry: company.OriginCountry,
		})
	}

	productionCountries := []model.ProductionCountry{}
	for _, country := range tmdbDetails.ProductionCountries {
		productionCountries = append(productionCountries, model.ProductionCountry{
			Iso3166_1: country.Iso3166_1,
			Name:      country.Name,
		})
	}

	return model.Movie{
		TmdbID:              tmdbDetails.ID,
		Title:               tmdbDetails.Title,
		Overview:            tmdbDetails.Overview,
		Tagline:             tmdbDetails.Tagline,
		Status:              tmdbDetails.Status,
		ReleaseDate:         tmdbDetails.ReleaseDate,
		Runtime:             tmdbDetails.Runtime,
		OriginalLanguage:    tmdbDetails.OriginalLanguage,
		OriginCountry:       tmdbDetails.OriginCountry,
		ProductionCompanies: productionCompanies,
		ProductionCountries: productionCountries,
		PosterPath:          tmdbDetails.PosterPath,
		BackdropPath:        tmdbDetails.BackdropPath,
	}, nil
}

func (s *MediaService) GetSeries(ctx context.Context, id int64) (model.Series, error) {
	tmdbDetails, err := s.tmdb.GetSeriesDetails(ctx, id)
	if err != nil {
		return model.Series{}, err
	}

	seasons := []model.Season{}
	for _, season := range tmdbDetails.Seasons {
		seasons = append(seasons, model.Season{
			TmdbID:       season.ID,
			SeasonNumber: int64(season.SeasonNumber),
			Title:        season.Name,
			Overview:     season.Overview,
			PosterPath:   season.PosterPath,
			AirDate:      season.AirDate,
		})
	}

	createdBy := []model.Person{}
	for _, person := range tmdbDetails.CreatedBy {
		createdBy = append(createdBy, model.Person{
			TmdbID:      person.ID,
			Name:        person.Name,
			ProfilePath: person.ProfilePath,
		})
	}

	networks := []model.Network{}
	for _, network := range tmdbDetails.Networks {
		networks = append(networks, model.Network{
			TmdbID:   network.ID,
			Name:     network.Name,
			LogoPath: network.LogoPath,
		})
	}

	genres := []model.Genre{}
	for _, genre := range tmdbDetails.Genres {
		genres = append(genres, model.Genre{
			TmdbID: genre.ID,
			Name:   genre.Name,
		})
	}

	return model.Series{
		TmdbID:       tmdbDetails.ID,
		Title:        tmdbDetails.Name,
		Overview:     tmdbDetails.Overview,
		Tagline:      tmdbDetails.Tagline,
		Status:       tmdbDetails.Status,
		Seasons:      seasons,
		PosterPath:   tmdbDetails.PosterPath,
		BackdropPath: tmdbDetails.BackdropPath,
		FirstAirDate: tmdbDetails.FirstAirDate,
		LastAirDate:  tmdbDetails.LastAirDate,
		InProduction: tmdbDetails.InProduction,
		CreatedBy:    createdBy,
		Networks:     networks,
		Genres:       genres,
	}, nil
}
