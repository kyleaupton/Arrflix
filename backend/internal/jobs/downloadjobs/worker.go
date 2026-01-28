package downloadjobs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	dbgen "github.com/kyleaupton/arrflix/internal/db/sqlc"
	"github.com/kyleaupton/arrflix/internal/downloader"
	"github.com/kyleaupton/arrflix/internal/importer"
	"github.com/kyleaupton/arrflix/internal/logger"
	"github.com/kyleaupton/arrflix/internal/repo"
	"github.com/kyleaupton/arrflix/internal/service"
	"github.com/kyleaupton/arrflix/internal/sse"
)

type Worker struct {
	repo     *repo.Repository
	dlm      *downloader.Manager
	importer *service.ImportService
	log      *logger.Logger
	broker   *sse.Broker

	pollInterval time.Duration
	claimLimit   int32
	maxAttempts  int
}

func New(repo *repo.Repository, dlm *downloader.Manager, importer *service.ImportService, log *logger.Logger, broker *sse.Broker) *Worker {
	return &Worker{
		repo:         repo,
		dlm:          dlm,
		importer:     importer,
		log:          log,
		broker:       broker,
		pollInterval: 3 * time.Second,
		claimLimit:   5,
		maxAttempts:  20,
	}
}

func (w *Worker) Run(ctx context.Context) {
	t := time.NewTicker(w.pollInterval)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			w.tick(ctx)
		}
	}
}

func (w *Worker) tick(ctx context.Context) {
	jobs, err := w.repo.ClaimRunnableDownloadJobs(ctx, w.claimLimit)
	if err != nil {
		w.log.Error().Err(err).Msg("download job claim failed")
		return
	}
	for _, job := range jobs {
		if err := w.processJob(ctx, job); err != nil {
			w.handleJobError(ctx, job, err)
		}
	}
}

func (w *Worker) processJob(ctx context.Context, job dbgen.DownloadJob) error {
	if job.Status == "cancelled" || job.Status == "failed" || job.Status == "imported" {
		return nil
	}

	client, err := w.dlm.GetClientByID(ctx, job.DownloaderID.String())
	if err != nil {
		return fmt.Errorf("get downloader client: %w", err)
	}

	if job.DownloaderExternalID == nil || *job.DownloaderExternalID == "" {
		addReq := downloader.AddRequest{}
		switch job.Protocol {
		case "torrent":
			addReq.MagnetURL = job.CandidateLink
		case "usenet":
			addReq.NZBURL = job.CandidateLink
		default:
			return fmt.Errorf("unknown protocol: %s", job.Protocol)
		}
		w.log.Info().Str("job_id", job.ID.String()).Str("add_req", fmt.Sprintf("%+v", addReq)).Msg("adding download")
		res, err := client.Add(ctx, addReq)
		if err != nil {
			return fmt.Errorf("downloader add: %w", err)
		}
		updated, err := w.repo.SetDownloadJobEnqueued(ctx, job.ID, res.ExternalID)
		if err == nil {
			w.publishJobUpdated(updated)
		}
		return err
	}

	externalID := *job.DownloaderExternalID
	item, err := client.Get(ctx, externalID)
	if err != nil {
		return fmt.Errorf("downloader get: %w", err)
	}

	w.log.Info().Interface("item", item).Msg("Item")

	jobStatus := mapItemStatus(item.Status)
	savePath := item.SavePath
	contentPath := item.ContentPath

	updated, err := w.repo.SetDownloadJobDownloadSnapshot(ctx, dbgen.SetDownloadJobDownloadSnapshotParams{
		ID:                  job.ID,
		Status:              jobStatus,
		DownloaderStatus:    ptr(string(item.Status)),
		Progress:            ptr(item.Progress),
		DownloadSavePath:    ptr(savePath),
		DownloadContentPath: ptr(contentPath),
	})
	if err != nil {
		return fmt.Errorf("update snapshot: %w", err)
	}
	w.publishJobUpdated(updated)

	// Terminal downloader error: mark job failed.
	if jobStatus == "failed" {
		failed, _ := w.repo.MarkDownloadJobFailed(ctx, job.ID, "downloader reported errored status")
		w.publishJobUpdated(failed)
		return nil
	}

	// Import when completed
	if jobStatus == "completed" && job.ImportDestPath == nil {
		if job.MediaType == "movie" {
			files, err := client.ListFiles(ctx, externalID)
			if err != nil && !errors.Is(err, downloader.ErrUnsupported) {
				return fmt.Errorf("list files: %w", err)
			}

			sourcePath := ""
			if len(files) > 0 {
				mainFile, ok := importer.PickMainMovieFile(files)
				if !ok {
					return fmt.Errorf("no files available for import")
				}
				// qBittorrent file paths are relative; use SavePath as root.
				if filepath.IsAbs(mainFile.Path) {
					sourcePath = mainFile.Path
				} else if item.SavePath != "" {
					sourcePath = filepath.Join(item.SavePath, mainFile.Path)
				} else if item.ContentPath != "" {
					// Best-effort fallback
					sourcePath = filepath.Join(item.ContentPath, mainFile.Path)
				}
			} else {
				// If file listing unsupported, try to use content path directly.
				sourcePath = item.ContentPath
			}

			if sourcePath == "" {
				return fmt.Errorf("unable to determine source path for import")
			}

			if importing, err := w.repo.SetDownloadJobImporting(ctx, job.ID, sourcePath); err != nil {
				return fmt.Errorf("mark importing: %w", err)
			} else {
				w.publishJobUpdated(importing)
			}

			res, err := w.importer.ImportMovieFile(ctx, job, sourcePath)
			if err != nil {
				return fmt.Errorf("import: %w", err)
			}

			method := res.Method
			imported, err := w.repo.SetDownloadJobImported(ctx, dbgen.SetDownloadJobImportedParams{
				ID:                 job.ID,
				ImportSourcePath:   &res.SourcePath,
				ImportDestPath:     &res.DestRelPath,
				ImportMethod:       &method,
				PrimaryMediaFileID: res.MediaFile.ID,
			})
			if err != nil {
				return fmt.Errorf("mark imported: %w", err)
			}
			if linkErr := w.repo.LinkDownloadJobMediaFile(ctx, job.ID, res.MediaFile.ID); linkErr != nil {
				w.log.Warn().Err(linkErr).Msg("failed to link download job to media file")
			}
			w.publishJobUpdated(imported)
		} else if job.MediaType == "series" {
			sourcePath := item.SavePath
			if sourcePath == "" {
				sourcePath = item.ContentPath
			}

			if importing, err := w.repo.SetDownloadJobImporting(ctx, job.ID, sourcePath); err != nil {
				return fmt.Errorf("mark importing: %w", err)
			} else {
				w.publishJobUpdated(importing)
			}

			results, err := w.importer.ImportSeriesJob(ctx, job, client)
			if err != nil {
				w.log.Error().Err(err).Str("job_id", job.ID.String()).Msg("ImportSeriesJob failed")
				return fmt.Errorf("import series: %w", err)
			}

			w.log.Debug().Int("results_count", len(results)).Str("job_id", job.ID.String()).Msg("ImportSeriesJob completed")

			if len(results) > 0 {
				// Mark job as imported using the first file as primary for the snapshot
				res := results[0]
				method := res.Method
				imported, err := w.repo.SetDownloadJobImported(ctx, dbgen.SetDownloadJobImportedParams{
					ID:                 job.ID,
					ImportSourcePath:   &sourcePath, // Use the root source path
					ImportDestPath:     &res.DestRelPath,
					ImportMethod:       &method,
					PrimaryMediaFileID: res.MediaFile.ID,
				})
				if err != nil {
					return fmt.Errorf("mark imported: %w", err)
				}
				w.publishJobUpdated(imported)
			} else {
				return fmt.Errorf("no episodes were successfully imported")
			}
		}
	}

	return nil
}

func (w *Worker) handleJobError(ctx context.Context, job dbgen.DownloadJob, err error) {
	msg := err.Error()
	w.log.Error().Err(err).Str("job_id", job.ID.String()).Msg("download job failed step")

	attempt := int(job.AttemptCount) + 1
	if attempt >= w.maxAttempts {
		failed, _ := w.repo.MarkDownloadJobFailed(ctx, job.ID, msg)
		w.publishJobUpdated(failed)
		return
	}

	backoff := time.Duration(attempt*attempt) * time.Second
	next := time.Now().Add(backoff)
	_, _ = w.repo.BumpDownloadJobRetry(ctx, dbgen.BumpDownloadJobRetryParams{
		LastError: &msg,
		NextRunAt: next,
		ID:        job.ID,
	})
}

func mapItemStatus(st downloader.JobStatus) string {
	switch st {
	case downloader.StatusCompleted, downloader.StatusSeeding:
		return "completed"
	case downloader.StatusQueued, downloader.StatusDownloading, downloader.StatusPaused:
		return "downloading"
	case downloader.StatusErrored:
		return "failed"
	default:
		return "downloading"
	}
}

func ptr[T any](v T) *T { return &v }

func (w *Worker) publishJobUpdated(job dbgen.DownloadJob) {
	if w.broker == nil {
		return
	}
	b, err := json.Marshal(job)
	if err != nil {
		return
	}
	w.broker.Publish(sse.Event{
		Type: "download_job_updated",
		ID:   job.ID.String(),
		Data: b,
	})
}
