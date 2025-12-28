package service

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	dbgen "github.com/kyleaupton/snaggle/backend/internal/db/sqlc"
	"github.com/kyleaupton/snaggle/backend/internal/downloader"
	"github.com/kyleaupton/snaggle/backend/internal/importer"
	"github.com/kyleaupton/snaggle/backend/internal/logger"
	"github.com/kyleaupton/snaggle/backend/internal/quality"
	"github.com/kyleaupton/snaggle/backend/internal/repo"
	"github.com/kyleaupton/snaggle/backend/internal/template"
)

type ImportService struct {
	repo *repo.Repository
	log  *logger.Logger
}

func NewImportService(r *repo.Repository, l *logger.Logger) *ImportService {
	return &ImportService{repo: r, log: l}
}

type ImportResult struct {
	SourcePath  string
	DestPath    string
	DestRelPath string
	Method      string // hardlink|copy
	MediaItem   dbgen.MediaItem
	MediaFile   dbgen.MediaFile
}

func (s *ImportService) ImportMovieFile(ctx context.Context, job dbgen.DownloadJob, sourcePath string) (ImportResult, error) {
	if !job.MediaItemID.Valid {
		return ImportResult{}, fmt.Errorf("job missing media_item_id")
	}

	mediaItem, err := s.repo.GetMediaItem(ctx, job.MediaItemID)
	if err != nil {
		return ImportResult{}, fmt.Errorf("get media item: %w", err)
	}

	year := ""
	if mediaItem.Year != nil {
		year = fmt.Sprintf("%d", *mediaItem.Year)
	}

	lib, err := s.repo.GetLibrary(ctx, job.LibraryID)
	if err != nil {
		return ImportResult{}, fmt.Errorf("get library: %w", err)
	}
	nt, err := s.repo.GetNameTemplate(ctx, job.NameTemplateID)
	if err != nil {
		return ImportResult{}, fmt.Errorf("get name template: %w", err)
	}

	// Parse quality from the candidate title if available
	parser := quality.NewParser()
	q := parser.Parse(job.CandidateTitle)

	context := quality.NamingContext{
		Title:   mediaItem.Title,
		Year:    year,
		Quality: q,
	}

	var rel string
	ext := filepath.Ext(sourcePath)

	// Use predicted_dest_path if available and replace .{ext} with actual extension
	if job.PredictedDestPath != nil && *job.PredictedDestPath != "" && strings.Contains(*job.PredictedDestPath, ".{ext}") {
		rel = strings.Replace(*job.PredictedDestPath, ".{ext}", ext, 1)
	} else {
		// Fallback to calculating from template (existing logic)
		if nt.Type == "series" {
			showPart, err := template.Render(coalesce(nt.SeriesShowTemplate, ""), context)
			if err != nil {
				return ImportResult{}, fmt.Errorf("render show template: %w", err)
			}
			seasonPart, err := template.Render(coalesce(nt.SeriesSeasonTemplate, ""), context)
			if err != nil {
				return ImportResult{}, fmt.Errorf("render season template: %w", err)
			}
			filePart, err := template.Render(nt.Template, context)
			if err != nil {
				return ImportResult{}, fmt.Errorf("render file template: %w", err)
			}
			rel = filepath.Join(showPart, seasonPart, filePart)
		} else {
			rel, err = template.Render(nt.Template, context)
			if err != nil {
				return ImportResult{}, fmt.Errorf("render template: %w", err)
			}
		}
		rel = importer.EnsureExt(rel, ext)
	}

	dest := filepath.Join(lib.RootPath, rel)

	destRel, err := filepath.Rel(lib.RootPath, dest)
	if err != nil || strings.HasPrefix(destRel, "..") {
		return ImportResult{}, fmt.Errorf("compute relative dest: %w", err)
	}

	method, err := importer.HardlinkOrCopy(sourcePath, dest)
	if err != nil {
		return ImportResult{}, err
	}

	// Upsert-ish media_file by path
	mf, err := s.repo.GetMediaFileByLibraryAndPath(ctx, lib.ID, destRel)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			mf, err = s.repo.CreateMediaFile(ctx, lib.ID, mediaItem.ID, (*pgtype.UUID)(nil), (*pgtype.UUID)(nil), destRel, nil)
			if err != nil {
				return ImportResult{}, fmt.Errorf("create media file: %w", err)
			}
		} else {
			return ImportResult{}, fmt.Errorf("get media file: %w", err)
		}
	}

	return ImportResult{
		SourcePath:  sourcePath,
		DestPath:    dest,
		DestRelPath: destRel,
		Method:      method,
		MediaItem:   mediaItem,
		MediaFile:   mf,
	}, nil
}

func (s *ImportService) ImportSeriesJob(ctx context.Context, job dbgen.DownloadJob, downloaderClient downloader.Client) ([]ImportResult, error) {
	if !job.MediaItemID.Valid {
		return nil, fmt.Errorf("job missing media_item_id")
	}

	mediaItem, err := s.repo.GetMediaItem(ctx, job.MediaItemID)
	if err != nil {
		return nil, fmt.Errorf("get media item: %w", err)
	}

	lib, err := s.repo.GetLibrary(ctx, job.LibraryID)
	if err != nil {
		return nil, fmt.Errorf("get library: %w", err)
	}
	nt, err := s.repo.GetNameTemplate(ctx, job.NameTemplateID)
	if err != nil {
		return nil, fmt.Errorf("get name template: %w", err)
	}

	if job.DownloaderExternalID == nil {
		return nil, fmt.Errorf("job missing downloader_external_id")
	}

	files, err := downloaderClient.ListFiles(ctx, *job.DownloaderExternalID)
	if err != nil {
		return nil, fmt.Errorf("list downloader files: %w", err)
	}

	var targetSeason *int
	if job.SeasonID.Valid {
		season, err := s.repo.GetSeason(ctx, job.SeasonID)
		if err == nil {
			sNum := int(season.SeasonNumber)
			targetSeason = &sNum
		}
	}

	var targetEpisode *int
	if job.EpisodeID.Valid {
		episode, err := s.repo.GetEpisode(ctx, job.EpisodeID)
		if err == nil {
			eNum := int(episode.EpisodeNumber)
			targetEpisode = &eNum
		}
	}

	matchedFiles := importer.MatchFilesToEpisodes(files, targetSeason, targetEpisode)
	s.log.Debug().
		Int("matched_count", len(matchedFiles)).
		Interface("matched_files", matchedFiles).
		Msg("Matched files to episodes")

	if len(matchedFiles) == 0 {
		return nil, fmt.Errorf("no files matched target episodes")
	}

	// Fetch downloader details to get SavePath
	dlItem, err := downloaderClient.Get(ctx, *job.DownloaderExternalID)
	if err != nil {
		return nil, fmt.Errorf("get downloader item: %w", err)
	}

	var results []ImportResult
	parser := quality.NewParser()
	q := parser.Parse(job.CandidateTitle)

	for epNum, f := range matchedFiles {
		s.log.Debug().Int("episode", epNum).Str("file", f.Path).Msg("Processing matched file")
		// Ensure we have a season for this episode
		var season dbgen.MediaSeason
		if targetSeason != nil && epNum >= 0 {
			season, err = s.repo.UpsertSeason(ctx, mediaItem.ID, int32(*targetSeason), pgtype.Date{Valid: false})
			if err != nil {
				s.log.Error().Err(err).Int("season", *targetSeason).Msg("Failed to upsert season")
				continue
			}
		} else {
			// If we don't know the season, we might need to parse it from the file
			info, ok := importer.ParseSeriesInfo(filepath.Base(f.Path))
			if !ok {
				s.log.Warn().Str("file", f.Path).Msg("Failed to parse series info from filename")
				continue
			}
			season, err = s.repo.UpsertSeason(ctx, mediaItem.ID, int32(info.Season), pgtype.Date{Valid: false})
			if err != nil {
				s.log.Error().Err(err).Int("season", info.Season).Msg("Failed to upsert parsed season")
				continue
			}
		}

		// Ensure we have an episode record
		episode, err := s.repo.UpsertEpisode(ctx, season.ID, int32(epNum), nil, pgtype.Date{Valid: false}, nil, nil)
		if err != nil {
			s.log.Error().Err(err).Int("episode", epNum).Msg("Failed to upsert episode")
			continue
		}

		epTitle := ""
		if episode.Title != nil {
			epTitle = *episode.Title
		}

		context := quality.NamingContext{
			Title:        mediaItem.Title,
			Year:         fmt.Sprintf("%d", *mediaItem.Year),
			Season:       fmt.Sprintf("%02d", season.SeasonNumber),
			Episode:      fmt.Sprintf("%02d", episode.EpisodeNumber),
			EpisodeTitle: epTitle,
			Quality:      q,
		}

		ext := filepath.Ext(f.Path)
		var rel string
		if nt.Type == "series" {
			showPart, err := template.Render(coalesce(nt.SeriesShowTemplate, ""), context)
			if err != nil {
				s.log.Error().Err(err).Interface("context", context).Msg("Failed to render show template")
				continue
			}
			seasonPart, err := template.Render(coalesce(nt.SeriesSeasonTemplate, ""), context)
			if err != nil {
				s.log.Error().Err(err).Interface("context", context).Msg("Failed to render season template")
				continue
			}
			filePart, err := template.Render(nt.Template, context)
			if err != nil {
				s.log.Error().Err(err).Interface("context", context).Msg("Failed to render file template")
				continue
			}
			rel = filepath.Join(showPart, seasonPart, filePart)
		} else {
			rel, err = template.Render(nt.Template, context)
			if err != nil {
				s.log.Error().Err(err).Interface("context", context).Msg("Failed to render template")
				continue
			}
		}
		rel = importer.EnsureExt(rel, ext)

		sourcePath := f.Path
		if !filepath.IsAbs(sourcePath) && dlItem.SavePath != "" {
			sourcePath = filepath.Join(dlItem.SavePath, f.Path)
		}

		dest := filepath.Join(lib.RootPath, rel)
		destRel, err := filepath.Rel(lib.RootPath, dest)
		if err != nil {
			s.log.Error().Err(err).Str("dest", dest).Msg("Failed to compute relative path")
			continue
		}

		s.log.Debug().Str("source", sourcePath).Str("dest", dest).Msg("Attempting import")
		method, err := importer.HardlinkOrCopy(sourcePath, dest)
		if err != nil {
			s.log.Error().Err(err).Str("source", sourcePath).Str("dest", dest).Msg("Failed to hardlink or copy file")
			continue
		}

		mf, err := s.repo.CreateMediaFile(ctx, lib.ID, mediaItem.ID, &season.ID, &episode.ID, destRel, nil)
		if err != nil {
			s.log.Error().Err(err).Str("path", destRel).Msg("Failed to create media file record")
			continue
		}

		if err := s.repo.LinkDownloadJobMediaFile(ctx, job.ID, mf.ID); err != nil {
			s.log.Warn().Err(err).Msg("Failed to link download job to media file")
		}

		s.log.Info().Str("path", destRel).Str("method", method).Msg("Successfully imported episode")
		results = append(results, ImportResult{
			SourcePath:  sourcePath,
			DestPath:    dest,
			DestRelPath: destRel,
			Method:      method,
			MediaItem:   mediaItem,
			MediaFile:   mf,
		})
	}

	return results, nil
}
