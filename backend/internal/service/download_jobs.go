package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	dbgen "github.com/kyleaupton/snaggle/backend/internal/db/sqlc"
	"github.com/kyleaupton/snaggle/backend/internal/repo"
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

func (s *DownloadJobsService) List(ctx context.Context) ([]dbgen.DownloadJob, error) {
	return s.repo.ListDownloadJobs(ctx)
}

func (s *DownloadJobsService) ListByMovie(ctx context.Context, tmdbMovieID int64) ([]dbgen.DownloadJob, error) {
	return s.repo.ListDownloadJobsByTmdbMovieID(ctx, tmdbMovieID)
}

func (s *DownloadJobsService) Cancel(ctx context.Context, id pgtype.UUID) (dbgen.DownloadJob, error) {
	job, err := s.repo.CancelDownloadJob(ctx, id)
	if err != nil {
		return dbgen.DownloadJob{}, fmt.Errorf("cancel job: %w", err)
	}
	return job, nil
}


