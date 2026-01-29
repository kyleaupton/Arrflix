// Package importw implements the import worker that processes import tasks
// (hardlinks/copies files from downloads to library).
package importw

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	dbgen "github.com/kyleaupton/arrflix/internal/db/sqlc"
	apperrors "github.com/kyleaupton/arrflix/internal/errors"
	"github.com/kyleaupton/arrflix/internal/importer"
	"github.com/kyleaupton/arrflix/internal/jobs/state"
	"github.com/kyleaupton/arrflix/internal/logger"
	"github.com/kyleaupton/arrflix/internal/model"
	"github.com/kyleaupton/arrflix/internal/release"
	"github.com/kyleaupton/arrflix/internal/repo"
	"github.com/kyleaupton/arrflix/internal/sse"
	"github.com/kyleaupton/arrflix/internal/template"
)

// Worker processes import tasks: hardlinks/copies files from downloads to library.
type Worker struct {
	repo   *repo.Repository
	log    *logger.Logger
	broker *sse.Broker
	sm     *state.ImportTaskMachine

	pollInterval time.Duration
	claimLimit   int32
	maxAttempts  int
}

// Config holds worker configuration.
type Config struct {
	PollInterval time.Duration
	ClaimLimit   int32
	MaxAttempts  int
}

// DefaultConfig returns default worker configuration.
func DefaultConfig() Config {
	return Config{
		PollInterval: 2 * time.Second,
		ClaimLimit:   10,
		MaxAttempts:  5,
	}
}

// New creates a new import worker.
func New(r *repo.Repository, log *logger.Logger, broker *sse.Broker) *Worker {
	cfg := DefaultConfig()
	return &Worker{
		repo:         r,
		log:          log,
		broker:       broker,
		sm:           state.NewImportTaskMachine(),
		pollInterval: cfg.PollInterval,
		claimLimit:   cfg.ClaimLimit,
		maxAttempts:  cfg.MaxAttempts,
	}
}

// Run starts the worker loop.
func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.pollInterval)
	defer ticker.Stop()

	w.log.Info().Msg("import worker started")

	for {
		select {
		case <-ctx.Done():
			w.log.Info().Msg("import worker stopped")
			return
		case <-ticker.C:
			w.tick(ctx)
		}
	}
}

func (w *Worker) tick(ctx context.Context) {
	// ClaimRunnableImportTasks atomically sets status to in_progress
	tasks, err := w.repo.ClaimRunnableImportTasks(ctx, w.claimLimit)
	if err != nil {
		w.log.Error().Err(err).Msg("failed to claim import tasks")
		return
	}

	for _, task := range tasks {
		if err := w.processTask(ctx, task); err != nil {
			w.handleError(ctx, task, err)
		}
	}
}

func (w *Worker) processTask(ctx context.Context, task dbgen.ImportTask) error {
	w.log.Info().
		Str("task_id", task.ID.String()).
		Str("source_path", task.SourcePath).
		Msg("processing import task")

	// Log status change to in_progress
	w.logEvent(ctx, task.ID, "status_changed", "", map[string]any{
		"old_status": "pending",
		"new_status": "in_progress",
	})

	// Validate source exists
	srcInfo, err := os.Stat(task.SourcePath)
	if err != nil {
		if os.IsNotExist(err) {
			return apperrors.AsPermanent(fmt.Errorf("source file not found: %s", task.SourcePath))
		}
		return fmt.Errorf("stat source: %w", err)
	}
	if srcInfo.IsDir() {
		return apperrors.AsPermanent(fmt.Errorf("source is a directory, expected file: %s", task.SourcePath))
	}

	// Get required data
	taskDetails, err := w.repo.GetImportTaskWithDetails(ctx, task.ID)
	if err != nil {
		return fmt.Errorf("get task details: %w", err)
	}

	// Compute destination path using name template
	destPath, err := w.computeDestPath(task, taskDetails)
	if err != nil {
		return apperrors.AsPermanent(fmt.Errorf("compute dest path: %w", err))
	}

	// Full absolute destination
	fullDest := filepath.Join(taskDetails.LibraryRootPath, destPath)

	// Check for existing destination
	if _, err := os.Stat(fullDest); err == nil {
		// Handle reimport: if this is a reimport, remove old file first
		if task.PreviousTaskID.Valid {
			w.log.Info().
				Str("task_id", task.ID.String()).
				Str("dest", fullDest).
				Msg("reimport: removing existing destination")
			if err := os.Remove(fullDest); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("remove existing dest: %w", err)
			}
		} else {
			// Not a reimport, fail
			return apperrors.AsPermanent(fmt.Errorf("destination already exists: %s", fullDest))
		}
	}

	// Perform import (hardlink or copy)
	method, err := importer.HardlinkOrCopy(task.SourcePath, fullDest)
	if err != nil {
		return fmt.Errorf("import file: %w", err)
	}

	w.log.Info().
		Str("task_id", task.ID.String()).
		Str("method", method).
		Str("dest", fullDest).
		Msg("file imported successfully")

	// Create media file record
	var episodeID *pgtype.UUID
	if task.EpisodeID.Valid {
		episodeID = &task.EpisodeID
	}

	mediaFile, err := w.repo.CreateMediaFile(ctx, task.LibraryID, task.MediaItemID, episodeID, destPath)
	if err != nil {
		// File was created but record failed - log but don't fail the task
		w.log.Error().Err(err).
			Str("task_id", task.ID.String()).
			Str("dest", destPath).
			Msg("failed to create media file record after successful import")
	}

	// Create file state
	if mediaFile.ID.Valid {
		fileSize := srcInfo.Size()
		_, _ = w.repo.UpsertMediaFileState(ctx, mediaFile.ID, true, &fileSize)
	}

	// Record import in media_file_import table
	_, _ = w.repo.CreateMediaFileImport(ctx, dbgen.CreateMediaFileImportParams{
		MediaFileID:  mediaFile.ID,
		ImportTaskID: task.ID,
		Method:       method,
		SourcePath:   &task.SourcePath,
		DestPath:     destPath,
		Success:      true,
		ErrorMessage: nil,
	})

	// Mark task completed
	_, err = w.repo.SetImportTaskCompleted(ctx, task.ID, destPath, method, mediaFile.ID)
	if err != nil {
		return fmt.Errorf("set task completed: %w", err)
	}

	w.logEvent(ctx, task.ID, "status_changed", "", map[string]any{
		"old_status":    "in_progress",
		"new_status":    "completed",
		"dest_path":     destPath,
		"import_method": method,
	})

	w.publishTaskUpdated(task.ID)
	return nil
}

func (w *Worker) computeDestPath(task dbgen.ImportTask, details dbgen.GetImportTaskWithDetailsRow) (string, error) {
	srcExt := filepath.Ext(task.SourcePath)

	// Build evaluation context for template rendering
	candidateTitle := ""
	if details.CandidateTitle != nil {
		candidateTitle = *details.CandidateTitle
	}

	q := release.Parse(candidateTitle)
	candidate := model.DownloadCandidate{
		Title: candidateTitle,
	}
	evalCtx := model.NewEvaluationContext(candidate, q)

	// Add media metadata
	year := 0
	if details.MediaYear != nil {
		year = int(*details.MediaYear)
	}
	tmdbID := int64(0) // Not available in details

	if task.MediaType == "movie" {
		evalCtx = evalCtx.WithMedia(model.MediaTypeMovie, details.MediaTitle, year, tmdbID)
	} else {
		evalCtx = evalCtx.WithMedia(model.MediaTypeSeries, details.MediaTitle, year, tmdbID)
		var seasonNum, epNum *int
		if details.SeasonNumber != nil {
			sn := int(*details.SeasonNumber)
			seasonNum = &sn
		}
		if details.EpisodeNumber != nil {
			en := int(*details.EpisodeNumber)
			epNum = &en
		}
		evalCtx = evalCtx.WithSeriesInfo(seasonNum, epNum, details.EpisodeTitle)
	}

	templateData := evalCtx.ToTemplateData()

	// Render template parts
	var rel string
	if task.MediaType == "series" {
		showPart, err := template.Render(coalesce(details.SeriesShowTemplate), templateData)
		if err != nil {
			return "", fmt.Errorf("render show template: %w", err)
		}
		seasonPart, err := template.Render(coalesce(details.SeriesSeasonTemplate), templateData)
		if err != nil {
			return "", fmt.Errorf("render season template: %w", err)
		}
		filePart, err := template.Render(details.NameTemplate, templateData)
		if err != nil {
			return "", fmt.Errorf("render file template: %w", err)
		}
		rel = filepath.Join(showPart, seasonPart, filePart)
	} else {
		// Movie type
		var dirPart string
		if details.MovieDirTemplate != nil && *details.MovieDirTemplate != "" {
			var err error
			dirPart, err = template.Render(*details.MovieDirTemplate, templateData)
			if err != nil {
				return "", fmt.Errorf("render movie dir template: %w", err)
			}
		}
		filePart, err := template.Render(details.NameTemplate, templateData)
		if err != nil {
			return "", fmt.Errorf("render file template: %w", err)
		}
		if dirPart != "" {
			rel = filepath.Join(dirPart, filePart)
		} else {
			rel = filePart
		}
	}

	rel = importer.EnsureExt(rel, srcExt)
	return rel, nil
}

func (w *Worker) handleError(ctx context.Context, task dbgen.ImportTask, err error) {
	msg := err.Error()
	category := apperrors.CategoryOf(err)

	w.log.Error().
		Err(err).
		Str("task_id", task.ID.String()).
		Str("category", string(category)).
		Msg("import task error")

	w.logEvent(ctx, task.ID, "error", msg, map[string]any{
		"category":      category,
		"attempt_count": task.AttemptCount + 1,
	})

	// Permanent errors fail immediately
	if category == apperrors.Permanent {
		_, _ = w.repo.SetImportTaskFailed(ctx, task.ID, msg, category)
		w.publishTaskUpdated(task.ID)
		return
	}

	// Check if we've exceeded max attempts
	attempt := int(task.AttemptCount) + 1
	maxAttempts := int(task.MaxAttempts)
	if maxAttempts == 0 {
		maxAttempts = w.maxAttempts
	}

	if attempt >= maxAttempts {
		_, _ = w.repo.SetImportTaskFailed(ctx, task.ID,
			fmt.Sprintf("max attempts (%d) exceeded: %s", maxAttempts, msg),
			apperrors.Transient)
		w.publishTaskUpdated(task.ID)
		return
	}

	// Schedule retry with exponential backoff
	backoff := time.Duration(math.Pow(2, float64(attempt))) * time.Second
	nextRun := time.Now().Add(backoff)

	w.logEvent(ctx, task.ID, "retry_scheduled", msg, map[string]any{
		"next_run_at": nextRun,
		"backoff":     backoff.String(),
	})

	_, _ = w.repo.ScheduleImportTaskRetry(ctx, task.ID, msg, category, nextRun)
	w.publishTaskUpdated(task.ID)
}

func (w *Worker) logEvent(ctx context.Context, taskID pgtype.UUID, eventType, message string, metadata map[string]any) {
	var metaBytes []byte
	if metadata != nil {
		metaBytes, _ = json.Marshal(metadata)
	}

	_, err := w.repo.CreateImportTaskEvent(ctx, dbgen.CreateImportTaskEventParams{
		ImportTaskID: taskID,
		EventType:    eventType,
		OldStatus:    nil,
		NewStatus:    nil,
		Message:      strPtr(message),
		Metadata:     metaBytes,
	})
	if err != nil {
		w.log.Warn().Err(err).Msg("failed to log import task event")
	}
}

func (w *Worker) publishTaskUpdated(taskID pgtype.UUID) {
	if w.broker == nil {
		return
	}
	w.broker.Publish(sse.Event{
		Type: "import_task_updated",
		ID:   taskID.String(),
		Data: nil,
	})
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func coalesce(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
