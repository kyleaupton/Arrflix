package service

import (
	"context"
	"strconv"

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

func (s *MediaService) GetMovieDetail(ctx context.Context, tmdbID int64) (model.MovieDetail, error) {
	tmdbDetails, err := s.tmdb.GetMovieDetails(ctx, tmdbID)
	if err != nil {
		return model.MovieDetail{}, err
	}

	var mediaItem dbgen.MediaItem
	local := true
	mediaItem, err = s.repo.GetMediaItemByTmdbIDAndType(ctx, tmdbID, string(model.MediaTypeMovie))
	if err != nil {
		local = false
	}

	var files []dbgen.ListMediaFilesForItemRow
	if local {
		files, _ = s.repo.ListMediaFilesForItem(ctx, mediaItem.ID)
	}

	fileInfos, availability := buildFileInfoAndAvailability(files)
	genres := make([]model.Genre, 0, len(tmdbDetails.Genres))
	for _, g := range tmdbDetails.Genres {
		genres = append(genres, model.Genre{TmdbID: g.ID, Name: g.Name})
	}
	year := parseYear(tmdbDetails.ReleaseDate)

	return model.MovieDetail{
		TmdbID:       tmdbDetails.ID,
		Title:        tmdbDetails.Title,
		Year:         year,
		Overview:     tmdbDetails.Overview,
		Tagline:      tmdbDetails.Tagline,
		Status:       tmdbDetails.Status,
		ReleaseDate:  tmdbDetails.ReleaseDate,
		Runtime:      tmdbDetails.Runtime,
		Genres:       genres,
		PosterPath:   tmdbDetails.PosterPath,
		BackdropPath: tmdbDetails.BackdropPath,
		Availability: availability,
		Files:        fileInfos,
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

func (s *MediaService) GetSeriesDetail(ctx context.Context, tmdbID int64) (model.SeriesDetail, error) {
	tmdbDetails, err := s.tmdb.GetSeriesDetails(ctx, tmdbID)
	if err != nil {
		return model.SeriesDetail{}, err
	}

	var mediaItem dbgen.MediaItem
	local := true
	mediaItem, err = s.repo.GetMediaItemByTmdbIDAndType(ctx, tmdbID, string(model.MediaTypeSeries))
	if err != nil {
		local = false
	}

	var files []dbgen.ListMediaFilesForItemRow
	var episodes []dbgen.ListEpisodeAvailabilityForSeriesRow
	if local {
		files, _ = s.repo.ListMediaFilesForItem(ctx, mediaItem.ID)
		episodes, _ = s.repo.ListEpisodeAvailabilityForSeries(ctx, mediaItem.ID)
	}

	fileInfos, availability := buildFileInfoAndAvailability(files)
	genres := make([]model.Genre, 0, len(tmdbDetails.Genres))
	for _, g := range tmdbDetails.Genres {
		genres = append(genres, model.Genre{TmdbID: g.ID, Name: g.Name})
	}
	seasons := buildSeasonDetails(episodes)
	year := parseYear(tmdbDetails.FirstAirDate)

	return model.SeriesDetail{
		TmdbID:       tmdbDetails.ID,
		Title:        tmdbDetails.Name,
		Year:         year,
		Overview:     tmdbDetails.Overview,
		Tagline:      tmdbDetails.Tagline,
		Status:       tmdbDetails.Status,
		FirstAirDate: tmdbDetails.FirstAirDate,
		LastAirDate:  tmdbDetails.LastAirDate,
		InProduction: tmdbDetails.InProduction,
		Genres:       genres,
		PosterPath:   tmdbDetails.PosterPath,
		BackdropPath: tmdbDetails.BackdropPath,
		Availability: availability,
		Files:        fileInfos,
		Seasons:      seasons,
	}, nil
}

func buildFileInfoAndAvailability(files []dbgen.ListMediaFilesForItemRow) ([]model.FileInfo, model.Availability) {
	fileInfos := make([]model.FileInfo, 0, len(files))
	libAgg := map[string]struct {
		count   int
		status  []string
	}{}

	for _, f := range files {
		libID := f.LibraryID.String()
		var seasonNum *int32
		if f.SeasonNumber != nil {
			seasonNum = f.SeasonNumber
		}
		var episodeNum *int32
		if f.EpisodeNumber != nil {
			episodeNum = f.EpisodeNumber
		}
		fileInfos = append(fileInfos, model.FileInfo{
			ID:            f.ID.String(),
			LibraryID:     libID,
			Path:          f.Path,
			Status:        f.Status,
			SeasonNumber:  seasonNum,
			EpisodeNumber: episodeNum,
		})
		entry := libAgg[libID]
		entry.count++
		entry.status = append(entry.status, f.Status)
		libAgg[libID] = entry
	}

	libraries := make([]model.LibraryAvailability, 0, len(libAgg))
	for libID, agg := range libAgg {
		libraries = append(libraries, model.LibraryAvailability{
			LibraryID:    libID,
			FileCount:    agg.count,
			StatusRollup: bestStatus(agg.status),
		})
	}

	return fileInfos, model.Availability{
		IsInLibrary: len(files) > 0,
		Libraries:   libraries,
	}
}

func buildSeasonDetails(rows []dbgen.ListEpisodeAvailabilityForSeriesRow) []model.SeasonDetail {
	if len(rows) == 0 {
		return nil
	}
	seasonsMap := map[int32][]model.EpisodeAvailability{}
	for _, r := range rows {
		seasonNum := r.SeasonNumber
		ep := model.EpisodeAvailability{
			SeasonNumber:  seasonNum,
			EpisodeNumber: r.EpisodeNumber,
			Title:         r.Title,
		}
		if r.AirDate.Valid {
			val := r.AirDate.Time.Format("2006-01-02")
			ep.AirDate = &val
		}
		if r.FileID.Valid {
			id := r.FileID.String()
			ep.FileID = &id
			ep.Available = true
		} else {
			ep.Available = false
		}
		seasonsMap[seasonNum] = append(seasonsMap[seasonNum], ep)
	}

	seasons := make([]model.SeasonDetail, 0, len(seasonsMap))
	for seasonNum, eps := range seasonsMap {
		seasons = append(seasons, model.SeasonDetail{
			SeasonNumber: seasonNum,
			Episodes:     eps,
		})
	}
	return seasons
}

func bestStatus(statuses []string) string {
	if len(statuses) == 0 {
		return "missing"
	}
	priority := map[string]int{
		"available":   5,
		"downloading": 4,
		"importing":   3,
		"missing":     2,
		"failed":      1,
		"deleted":     0,
	}
	best := "missing"
	bestScore := -1
	for _, s := range statuses {
		if score, ok := priority[s]; ok {
			if score > bestScore {
				best = s
				bestScore = score
			}
		}
	}
	return best
}

func parseYear(dateStr string) *int32 {
	if len(dateStr) < 4 {
		return nil
	}
	y, err := strconv.Atoi(dateStr[:4])
	if err != nil {
		return nil
	}
	val := int32(y)
	return &val
}
