package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	dbgen "github.com/kyleaupton/arrflix/internal/db/sqlc"
	"github.com/kyleaupton/arrflix/internal/repo"
)

type DownloadJobsService struct {
	repo *repo.Repository
}

func NewDownloadJobsService(r *repo.Repository) *DownloadJobsService {
	return &DownloadJobsService{repo: r}
}

func (s *DownloadJobsService) Create(ctx context.Context, arg dbgen.CreateDownloadJobParams) (dbgen.DownloadJob, error) {
	return s.repo.CreateDownloadJob(ctx, arg)
}

func (s *DownloadJobsService) Get(ctx context.Context, id pgtype.UUID) (dbgen.DownloadJob, error) {
	return s.repo.GetDownloadJob(ctx, id)
}

// GetWithImportSummary returns a download job with computed import status summary.
func (s *DownloadJobsService) GetWithImportSummary(ctx context.Context, id pgtype.UUID) (dbgen.GetDownloadJobWithImportSummaryRow, error) {
	return s.repo.GetDownloadJobWithImportSummary(ctx, id)
}

// GetTimeline returns the combined event log for a download job (download events + import events).
func (s *DownloadJobsService) GetTimeline(ctx context.Context, id pgtype.UUID) ([]dbgen.GetDownloadJobTimelineRow, error) {
	return s.repo.GetDownloadJobTimeline(ctx, id)
}

func (s *DownloadJobsService) List(ctx context.Context) ([]dbgen.DownloadJob, error) {
	return s.repo.ListDownloadJobs(ctx)
}

func (s *DownloadJobsService) ListByMovie(ctx context.Context, tmdbMovieID int64) ([]dbgen.DownloadJob, error) {
	return s.repo.ListDownloadJobsByTmdbMovieID(ctx, tmdbMovieID)
}

func (s *DownloadJobsService) ListBySeries(ctx context.Context, tmdbSeriesID int64) ([]dbgen.ListDownloadJobsByTmdbSeriesIDRow, error) {
	return s.repo.ListDownloadJobsByTmdbSeriesID(ctx, tmdbSeriesID)
}

// Cancel cancels a download job and all its pending import tasks.
func (s *DownloadJobsService) Cancel(ctx context.Context, id pgtype.UUID) (dbgen.DownloadJob, error) {
	job, err := s.repo.CancelDownloadJob(ctx, id)
	if err != nil {
		return dbgen.DownloadJob{}, fmt.Errorf("cancel job: %w", err)
	}

	// Also cancel all pending import tasks for this job
	if err := s.repo.CancelPendingImportTasksForJob(ctx, id); err != nil {
		// Log but don't fail the cancel operation
		_ = err
	}

	return job, nil
}

// ListImportTasks returns all import tasks for a download job.
func (s *DownloadJobsService) ListImportTasks(ctx context.Context, jobID pgtype.UUID) ([]dbgen.ImportTask, error) {
	return s.repo.ListImportTasksByDownloadJob(ctx, jobID)
}
