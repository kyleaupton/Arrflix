package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	dbgen "github.com/kyleaupton/arrflix/internal/db/sqlc"
)

type DownloadJobsRepo interface {
	CreateDownloadJob(ctx context.Context, arg dbgen.CreateDownloadJobParams) (dbgen.DownloadJob, error)
	GetDownloadJob(ctx context.Context, id pgtype.UUID) (dbgen.DownloadJob, error)
	GetDownloadJobByCandidate(ctx context.Context, indexerID int64, guid string) (dbgen.DownloadJob, error)
	ListDownloadJobsByMediaItem(ctx context.Context, mediaItemID pgtype.UUID) ([]dbgen.DownloadJob, error)
	ListDownloadJobsByTmdbMovieID(ctx context.Context, tmdbMovieID int64) ([]dbgen.DownloadJob, error)
	ListDownloadJobsByTmdbSeriesID(ctx context.Context, tmdbSeriesID int64) ([]dbgen.ListDownloadJobsByTmdbSeriesIDRow, error)
	ListDownloadJobs(ctx context.Context) ([]dbgen.DownloadJob, error)
	CancelDownloadJob(ctx context.Context, id pgtype.UUID) (dbgen.DownloadJob, error)

	ClaimRunnableDownloadJobs(ctx context.Context, limit int32) ([]dbgen.DownloadJob, error)

	SetDownloadJobEnqueued(ctx context.Context, id pgtype.UUID, downloaderExternalID string) (dbgen.DownloadJob, error)
	SetDownloadJobDownloadSnapshot(ctx context.Context, arg dbgen.SetDownloadJobDownloadSnapshotParams) (dbgen.DownloadJob, error)
	SetDownloadJobImporting(ctx context.Context, id pgtype.UUID, importSourcePath string) (dbgen.DownloadJob, error)
	SetDownloadJobImported(ctx context.Context, arg dbgen.SetDownloadJobImportedParams) (dbgen.DownloadJob, error)
	LinkDownloadJobMediaFile(ctx context.Context, downloadJobID, mediaFileID pgtype.UUID) error

	BumpDownloadJobRetry(ctx context.Context, arg dbgen.BumpDownloadJobRetryParams) (dbgen.DownloadJob, error)
	MarkDownloadJobFailed(ctx context.Context, id pgtype.UUID, lastError string) (dbgen.DownloadJob, error)
}

func (r *Repository) CreateDownloadJob(ctx context.Context, arg dbgen.CreateDownloadJobParams) (dbgen.DownloadJob, error) {
	return r.Q.CreateDownloadJob(ctx, arg)
}

func (r *Repository) GetDownloadJob(ctx context.Context, id pgtype.UUID) (dbgen.DownloadJob, error) {
	return r.Q.GetDownloadJob(ctx, id)
}

func (r *Repository) GetDownloadJobByCandidate(ctx context.Context, indexerID int64, guid string) (dbgen.DownloadJob, error) {
	return r.Q.GetDownloadJobByCandidate(ctx, dbgen.GetDownloadJobByCandidateParams{
		IndexerID: indexerID,
		Guid:      guid,
	})
}

func (r *Repository) ListDownloadJobsByMediaItem(ctx context.Context, mediaItemID pgtype.UUID) ([]dbgen.DownloadJob, error) {
	return r.Q.ListDownloadJobsByMediaItem(ctx, mediaItemID)
}

func (r *Repository) ListDownloadJobsByTmdbMovieID(ctx context.Context, tmdbMovieID int64) ([]dbgen.DownloadJob, error) {
	return r.Q.ListDownloadJobsByTmdbMovieID(ctx, &tmdbMovieID)
}

func (r *Repository) ListDownloadJobsByTmdbSeriesID(ctx context.Context, tmdbSeriesID int64) ([]dbgen.ListDownloadJobsByTmdbSeriesIDRow, error) {
	return r.Q.ListDownloadJobsByTmdbSeriesID(ctx, &tmdbSeriesID)
}

func (r *Repository) ListDownloadJobs(ctx context.Context) ([]dbgen.DownloadJob, error) {
	return r.Q.ListDownloadJobs(ctx)
}

func (r *Repository) CancelDownloadJob(ctx context.Context, id pgtype.UUID) (dbgen.DownloadJob, error) {
	return r.Q.CancelDownloadJob(ctx, id)
}

func (r *Repository) ClaimRunnableDownloadJobs(ctx context.Context, limit int32) ([]dbgen.DownloadJob, error) {
	return r.Q.ClaimRunnableDownloadJobs(ctx, limit)
}

func (r *Repository) SetDownloadJobEnqueued(ctx context.Context, id pgtype.UUID, downloaderExternalID string) (dbgen.DownloadJob, error) {
	return r.Q.SetDownloadJobEnqueued(ctx, dbgen.SetDownloadJobEnqueuedParams{
		ID:                   id,
		DownloaderExternalID: &downloaderExternalID,
	})
}

func (r *Repository) SetDownloadJobDownloadSnapshot(ctx context.Context, arg dbgen.SetDownloadJobDownloadSnapshotParams) (dbgen.DownloadJob, error) {
	return r.Q.SetDownloadJobDownloadSnapshot(ctx, arg)
}

func (r *Repository) SetDownloadJobImporting(ctx context.Context, id pgtype.UUID, importSourcePath string) (dbgen.DownloadJob, error) {
	return r.Q.SetDownloadJobImporting(ctx, dbgen.SetDownloadJobImportingParams{
		ID:               id,
		ImportSourcePath: &importSourcePath,
	})
}

func (r *Repository) SetDownloadJobImported(ctx context.Context, arg dbgen.SetDownloadJobImportedParams) (dbgen.DownloadJob, error) {
	return r.Q.SetDownloadJobImported(ctx, arg)
}

func (r *Repository) LinkDownloadJobMediaFile(ctx context.Context, downloadJobID, mediaFileID pgtype.UUID) error {
	return r.Q.LinkDownloadJobMediaFile(ctx, dbgen.LinkDownloadJobMediaFileParams{
		DownloadJobID: downloadJobID,
		MediaFileID:   mediaFileID,
	})
}

func (r *Repository) BumpDownloadJobRetry(ctx context.Context, arg dbgen.BumpDownloadJobRetryParams) (dbgen.DownloadJob, error) {
	return r.Q.BumpDownloadJobRetry(ctx, arg)
}

func (r *Repository) MarkDownloadJobFailed(ctx context.Context, id pgtype.UUID, lastError string) (dbgen.DownloadJob, error) {
	return r.Q.MarkDownloadJobFailed(ctx, dbgen.MarkDownloadJobFailedParams{
		ID:        id,
		LastError: &lastError,
	})
}
