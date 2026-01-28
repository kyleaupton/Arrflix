package service

import (
	"context"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kyleaupton/arrflix/internal/identity"
	"github.com/kyleaupton/arrflix/internal/logger"
	"github.com/kyleaupton/arrflix/internal/repo"
)

type ScannerService struct {
	repo   *repo.Repository
	logger *logger.Logger
	tmdb   *TmdbService
}

func NewScannerService(r *repo.Repository, l *logger.Logger, tmdb *TmdbService) *ScannerService {
	return &ScannerService{repo: r, logger: l, tmdb: tmdb}
}

type ScanStats struct {
	FilesSeen         int `json:"filesSeen"`
	MediaItemsCreated int `json:"mediaItemsCreated"`
	Duration          int `json:"duration"`
}

func (s *ScannerService) StartScan(ctx context.Context, libraryID pgtype.UUID) (ScanStats, error) {
	stats := ScanStats{}
	start := time.Now()

	library, err := s.repo.GetLibrary(ctx, libraryID)
	if err != nil {
		s.logger.Error().Str("library_id", libraryID.String()).Err(err).Msg("Error getting library")
		return ScanStats{}, err
	}

	s.logger.Info().Str("library_name", library.Name).Str("library_path", library.RootPath).Msg("Starting Scan")

	err = filepath.WalkDir(library.RootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err // propagate permission/IO errors
		}

		if d.IsDir() || !isMediaFile(path) {
			s.logger.Debug().Str("path", path).Msg("Skipping Directory or Non-Media File")
			return nil
		}

		s.logger.Debug().Str("path", path).Msg("Processing Media File")

		stats.FilesSeen++

		relPath, err := filepath.Rel(library.RootPath, path)
		if err != nil || strings.HasPrefix(relPath, "..") {
			s.logger.Error().Str("path", path).Err(err).Msg("Path outside library root")
			return nil
		}

		// see if path exists in media_file
		// if it does, skip
		_, err = s.repo.GetMediaFileByLibraryAndPath(ctx, library.ID, relPath)
		if err != nil && err != pgx.ErrNoRows {
			return err
		}
		if err == nil {
			s.logger.Debug().Str("path", path).Msg("Media file already exists")
			return nil
		}

		// attempt to get identification
		identity, err := identity.Resolve(library, path)
		if err != nil {
			s.logger.Error().Str("path", path).Err(err).Msg("Error resolving identity")
			return nil
		}

		if identity.TmdbID == nil {
			// If we got an identity, but no tmdb id, we'll need to convert to one.
			// This is a best effort to get a tmdb id. If we don't get one, we'll skip the file.
			var id string
			var provider string

			if identity.TvdbID != nil {
				id = strconv.FormatInt(*identity.TvdbID, 10)
				provider = "tvdb_id"
			} else if identity.ImdbID != nil {
				id = *identity.ImdbID
				provider = "imdb_id"
			}

			s.logger.Debug().Str("path", path).Str("provider", provider).Str("id", id).Msg("Identity has no tmdb id, converting to tmdb id")

			// This is like a general "search for this external id" type of thing.
			// It will return a list of results categorized by the type of media it is.
			// For now, we'll just grab the first item from the right category and use that.
			// TODO: We should probably do something more intelligent here.
			res, err := s.tmdb.FindByID(ctx, id, provider)
			if err != nil {
				s.logger.Error().Str("path", path).Err(err).Msg("Error getting find by id")
				return nil
			}

			var tmdbId *int64
			switch library.Type {
			case "movie":
				if len(res.MovieResults) == 0 {
					s.logger.Error().Str("path", path).Msg("No movie results found")
					return nil
				}

				tmdbId = &res.MovieResults[0].ID
			case "series":
				if len(res.TvResults) == 0 {
					s.logger.Error().Str("path", path).Msg("No series results found")
					return nil
				}

				tmdbId = &res.TvResults[0].ID
			}

			identity.TmdbID = tmdbId
		}

		if identity.TmdbID == nil {
			// We can't do anything without a tmdb id, so we'll skip the file.
			// TODO: keep track of files that we couldn't do anything with so we can alert the user.
			s.logger.Error().Str("path", path).Msg("No tmdb id found")
			return nil
		}

		// See if the tmdb id exists within media_item
		mediaItem, err := s.repo.GetMediaItemByTmdbID(ctx, *identity.TmdbID)
		if err != nil && err != pgx.ErrNoRows {
			s.logger.Error().Str("path", path).Err(err).Msg("Error getting media item by tmdb id")
			return err
		}

		var mediaItemId pgtype.UUID
		var seasonId *pgtype.UUID
		var episodeId *pgtype.UUID

		if err == pgx.ErrNoRows {
			// If we get here then mediaItem doesn't exists, so we need to make it.
			s.logger.Debug().Str("path", path).Msg("Media item not found, grabbing things")

			var title string
			var year int32

			switch library.Type {
			case "movie":
				movie, err := s.tmdb.GetMovieDetails(ctx, *identity.TmdbID)
				if err != nil {
					s.logger.Error().Str("path", path).Err(err).Msg("Error getting movie details")
					return nil
				}

				title = movie.Title
				// yyyy-mm-dd
				yearStr := strings.Split(movie.ReleaseDate, "-")[0]
				// cast to int
				year64, err := strconv.ParseInt(yearStr, 10, 32)
				if err != nil {
					s.logger.Error().Str("path", path).Err(err).Msg("Error getting movie details")
					return nil
				}

				year = int32(year64)
			case "series":
				tv, err := s.tmdb.GetSeriesDetails(ctx, *identity.TmdbID)
				if err != nil {
					s.logger.Error().Str("path", path).Err(err).Msg("Error getting series details")
					return nil
				}

				title = tv.Name
				// yyyy-mm-dd
				yearStr := strings.Split(tv.FirstAirDate, "-")[0]
				// cast to int
				year64, err := strconv.ParseInt(yearStr, 10, 32)
				if err != nil {
					s.logger.Error().Str("path", path).Err(err).Msg("Error getting episode details")
					return nil
				}

				year = int32(year64)
			}

			s.logger.Debug().Str("title", title).Int32("year", year).Int64("tmdb_id", *identity.TmdbID).Msg("Creating media_item")

			// create media_item if it doesn't exist
			createdMediaItem, err := s.repo.CreateMediaItem(
				ctx,
				library.Type,
				title,
				&year,
				identity.TmdbID,
			)
			if err != nil {
				s.logger.Error().Str("path", path).Err(err).Msg("Error creating media_item")
				return nil
			}

			s.logger.Debug().Str("path", path).Str("media_item_id", createdMediaItem.ID.String()).Msg("Media item created")
			mediaItemId = createdMediaItem.ID
			stats.MediaItemsCreated++

			if identity.Season != nil {
				// create media_season if it doesn't exist
				seasonRow, err := s.repo.UpsertSeason(ctx, mediaItemId, *identity.Season, pgtype.Date{})
				if err != nil {
					return nil
				}
				seasonId = &seasonRow.ID

				if identity.Episode != nil {
					// grab the episode details from tmdb
					episode, err := s.tmdb.GetEpisodeDetails(ctx, *identity.TmdbID, int64(*identity.Season), int64(*identity.Episode))
					if err != nil {
						return err
					}

					// create media_episode if it doesn't exist
					episodeRow, err := s.repo.UpsertEpisode(ctx, seasonRow.ID, *identity.Episode, &episode.Name, pgtype.Date{}, nil, nil)
					if err != nil {
						return nil
					}
					episodeId = &episodeRow.ID
				}
			}

		} else {
			mediaItemId = mediaItem.ID
		}

		// create media_file
		s.repo.CreateMediaFile(ctx, library.ID, mediaItemId, seasonId, episodeId, relPath, nil)

		return nil
	})

	if err != nil {
		s.logger.Error().Err(err).Msg("Error walking directory")
		return ScanStats{}, err
	}

	s.logger.Info().Str("library_id", libraryID.String()).Msg("Scan Complete")
	stats.Duration = int(time.Since(start).Seconds())

	return stats, nil
}

func isMediaFile(path string) bool {
	extensions := []string{".mkv", ".mp4", ".avi", ".mov", ".wmv", ".flv", ".m4v", ".webm"}
	for _, ext := range extensions {
		if strings.HasSuffix(path, ext) {
			return true
		}
	}
	return false
}
