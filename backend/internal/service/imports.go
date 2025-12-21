package service

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	dbgen "github.com/kyleaupton/snaggle/backend/internal/db/sqlc"
	"github.com/kyleaupton/snaggle/backend/internal/importer"
	"github.com/kyleaupton/snaggle/backend/internal/quality"
	"github.com/kyleaupton/snaggle/backend/internal/repo"
	"github.com/kyleaupton/snaggle/backend/internal/template"
)

type ImportService struct {
	repo *repo.Repository
}

func NewImportService(r *repo.Repository) *ImportService {
	return &ImportService{repo: r}
}

type ImportResult struct {
	SourcePath string
	DestPath   string
	Method     string // hardlink|copy
	MediaItem  dbgen.MediaItem
	MediaFile  dbgen.MediaFile
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

	rel, err := template.Render(nt.Template, context)
	if err != nil {
		return ImportResult{}, fmt.Errorf("render template: %w", err)
	}

	ext := filepath.Ext(sourcePath)
	dest := filepath.Join(lib.RootPath, rel)
	dest = importer.EnsureExt(dest, ext)

	method, err := importer.HardlinkOrCopy(sourcePath, dest)
	if err != nil {
		return ImportResult{}, err
	}

	// Upsert-ish media_file by path
	mf, err := s.repo.GetMediaFileByPath(ctx, dest)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			mf, err = s.repo.CreateMediaFile(ctx, mediaItem.ID, (*pgtype.UUID)(nil), (*pgtype.UUID)(nil), dest)
			if err != nil {
				return ImportResult{}, fmt.Errorf("create media file: %w", err)
			}
		} else {
			return ImportResult{}, fmt.Errorf("get media file: %w", err)
		}
	}

	return ImportResult{
		SourcePath: sourcePath,
		DestPath:   dest,
		Method:     method,
		MediaItem:  mediaItem,
		MediaFile:  mf,
	}, nil
}
