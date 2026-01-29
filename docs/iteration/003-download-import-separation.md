# Iteration 003: Download/Import System Separation

**Status**: In Review
**Started**: 2026-01-29
**Goal**: Separate downloading from importing into distinct subsystems with explicit state machines, error classification, and immutable audit trails

---

## Problem Statement

The download system combined downloading and importing into a single monolithic flow with multiple issues:

1. **Combined concerns** - `download_job` tracked both download states (created, downloading) AND import states (importing, import_failed, import_complete) in one entity
2. **Single file assumption** - Schema assumed one file per download, but season packs have many episodes
3. **No retry granularity** - If 1 of 10 episodes failed import, the entire job was marked failed
4. **No reimport capability** - Can't retry a failed import without re-downloading
5. **Implicit state transitions** - No validation that state changes were valid (e.g., "created" → "completed" directly)
6. **No error classification** - All errors treated equally; couldn't distinguish permanent failures (invalid magnet) from transient ones (network timeout)
7. **No audit trail** - No visibility into what happened during download/import lifecycle

---

## Design Decisions

### 1. Separate Download Jobs from Import Tasks

**Before:**
```
download_job
├── id, status (8 states including import states)
├── download fields (protocol, link, external_id, progress)
├── import fields (dest_path, import_method)
└── media references (media_item_id, season_id, episode_id)
```

**After:**
```
download_job                    (6 states, download only)
├── id, status (created, enqueued, downloading, completed, failed, cancelled)
├── download fields (protocol, link, external_id, progress, save_path, content_path)
├── retry logic (attempt_count, next_run_at, last_error, error_category)
└── media references

import_task                     (5 states, per-file import)
├── id, status (pending, in_progress, completed, failed, cancelled)
├── download_job_id             (nullable FK)
├── previous_task_id            (reimport chain)
├── source_path, dest_path
├── retry logic (attempt_count, max_attempts, next_run_at, last_error, error_category)
└── media references
```

### 2. One Import Task Per File

Season packs now spawn multiple `import_task` rows:
- 10-episode pack → 10 import tasks
- Each tracks independently: 8 succeed, 2 fail → `partial_failure` status
- Individual reimport without re-downloading

### 3. Explicit State Machines

State transitions are validated before execution:

```go
var DownloadJobTransitions = map[Status][]Status{
    "created":     {"enqueued", "failed", "cancelled"},
    "enqueued":    {"downloading", "failed", "cancelled"},
    "downloading": {"completed", "failed", "cancelled"},
    "completed":   {},  // terminal
    "failed":      {},  // terminal
    "cancelled":   {},  // terminal
}

var ImportTaskTransitions = map[Status][]Status{
    "pending":     {"in_progress", "cancelled"},
    "in_progress": {"completed", "failed"},
    "completed":   {},  // terminal - reimport creates new task
    "failed":      {},  // terminal - reimport creates new task
    "cancelled":   {},  // terminal
}
```

### 4. Error Categories

Errors are classified for intelligent retry behavior:

```go
type Category string
const (
    Transient Category = "transient"  // retry: network timeout, 5xx
    Permanent Category = "permanent"  // fail fast: 401/403, invalid magnet
)
```

Default: unknown errors → transient (allow retry)

### 5. Event Tables for Audit Trail

Every state change and error is logged:

```sql
download_job_event  (created, status_changed, error, retry_scheduled)
import_task_event   (created, status_changed, error, retry_scheduled, reimport_requested)
```

Combined timeline query aggregates both for unified history view.

### 6. Computed Import Status

`GetDownloadJobWithImportSummary` returns aggregate status:
- `download_pending` - download not yet completed
- `awaiting_import` - completed, no import tasks yet
- `importing` - tasks in progress
- `fully_imported` - all tasks completed
- `partial_failure` - some completed, some failed
- `import_failed` - all tasks failed

### 7. Cancel Cascade

Cancelling a download job automatically cancels all pending import tasks for that job.

### 8. Reimport Chain

Import tasks have `previous_task_id` FK:
- Original import: `previous_task_id = NULL`
- Reimport: `previous_task_id = <original task ID>`
- Recursive CTE retrieves full chain history

Reimport handling: if destination exists and this is a reimport, remove old file first.

---

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│ internal/jobs/download/worker.go                            │
│                                                             │
│   Polls download clients, manages download_job lifecycle    │
│   Config: 3s interval, 20 claim limit, 10 max attempts      │
│                                                             │
│   created → enqueue to client → enqueued                    │
│   enqueued/downloading → poll status → downloading/completed│
│   completed → spawn import_task(s)                          │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│ internal/jobs/import/worker.go                              │
│                                                             │
│   Processes import tasks: hardlink/copy to library          │
│   Config: 2s interval, 10 claim limit, 5 max attempts       │
│                                                             │
│   pending → claim (atomic) → in_progress                    │
│   validate source → compute dest → hardlink/copy → completed│
└─────────────────────────────────────────────────────────────┘
```

---

## Schema Changes

### New Tables

| Table | Purpose |
|-------|---------|
| `import_task` | Per-file import tracking with reimport chain |
| `download_job_event` | Download job audit log |
| `import_task_event` | Import task audit log |

### Modified Tables

| Table | Change |
|-------|--------|
| `download_job` | Simplified to 6 states, added error_category, removed import fields |
| `media_file_import` | Added `import_task_id` FK |

### Deleted Tables/Columns

| Change | Reason |
|--------|--------|
| `download_job_media_file` | Replaced by `import_task` |
| `download_job.import_*` columns | Moved to `import_task` |

---

## New API Endpoints

```
# Import Tasks
GET    /v1/import-tasks              - List import tasks (paginated, filterable by status)
GET    /v1/import-tasks/counts       - Get counts by status
GET    /v1/import-tasks/:id          - Get import task with details
GET    /v1/import-tasks/:id/timeline - Get event log
GET    /v1/import-tasks/:id/history  - Get reimport chain
POST   /v1/import-tasks/:id/reimport - Create new task for reimport
POST   /v1/import-tasks/:id/cancel   - Cancel pending task

# Download Jobs (enhanced)
GET    /v1/download-jobs/:id         - Now returns import summary
GET    /v1/download-jobs/:id/timeline - Combined download + import events
GET    /v1/download-jobs/:id/import-tasks - List import tasks for job
```

---

## Files Created

| File | Purpose |
|------|---------|
| `internal/errors/category.go` | Error categorization (transient/permanent) |
| `internal/jobs/state/machine.go` | State machine definitions and validation |
| `internal/jobs/download/worker.go` | Download worker (replaces old worker) |
| `internal/jobs/import/worker.go` | Import worker |
| `internal/service/import_tasks.go` | Import tasks service |
| `internal/http/handlers/import_tasks.go` | Import tasks HTTP handler |
| `internal/repo/import_tasks.go` | Import tasks repository |
| `internal/db/queries/import_tasks.sql` | Import task SQLC queries |
| `internal/db/queries/download_job_events.sql` | Download event SQLC queries |
| `internal/db/queries/import_task_events.sql` | Import event SQLC queries |

---

## Files Modified

| File | Change |
|------|--------|
| `internal/db/migrations/0008_download_job.up.sql` | Complete rewrite with new schema |
| `internal/db/queries/download_jobs.sql` | Simplified states, added timeline/summary queries |
| `internal/service/download_jobs.go` | Added cancel cascade, timeline, import summary |
| `internal/http/handlers/download_jobs.go` | Added timeline, import-tasks endpoints |
| `internal/http/http.go` | Register import tasks handler |
| `cmd/api/main.go` | Start both download and import workers |

---

## Files Deleted

| File | Reason |
|------|--------|
| `internal/jobs/downloadjobs/worker.go` | Replaced by separate download/import workers |

---

## Worker Configuration

| Worker | Poll Interval | Claim Limit | Max Attempts |
|--------|---------------|-------------|--------------|
| Download | 3s | 20 | 10 |
| Import | 2s | 10 | 5 |

Both use exponential backoff: 2^attempt seconds

---

## Testing Checklist

- [ ] Nuke DB, verify migrations run clean
- [ ] `sqlc generate` succeeds
- [ ] `go build ./...` succeeds
- [ ] Search → enqueue download → verify `download_job` created with status `created`
- [ ] Watch status progress: `created` → `enqueued` → `downloading` → `completed`
- [ ] Verify `import_task`(s) spawned when download completes
- [ ] Verify import progresses: `pending` → `in_progress` → `completed`
- [ ] Check `GET /v1/download-jobs/:id` returns `import_status: fully_imported`
- [ ] Test reimport: `POST /v1/import-tasks/:id/reimport`
- [ ] Verify new task created with `previous_task_id` set
- [ ] Test timeline: `GET /v1/download-jobs/:id/timeline`
- [ ] Test cancel cascade: cancel download, verify pending imports cancelled
- [ ] Test series pack: download season, verify multiple import tasks created
- [ ] Simulate partial failure: verify `import_status: partial_failure`
- [ ] Test permanent error (stop qBittorrent): verify job fails after retries
- [ ] Test transient error: verify retry with exponential backoff

---

## Known Issues / TODOs

1. **qBittorrent error categorization not implemented** - Wrapper should mark 401/403 as permanent errors
2. **No "created" event logged for spawned import tasks** - Download worker creates tasks but doesn't log events
3. **ListDownloadJobs has no pagination** - Could be slow with many jobs

---

## Future Work (Out of Scope)

- Import task cleanup/retention policy (delete old completed tasks)
- Frontend UI for import task management
- Configurable worker limits via settings
- Metrics/monitoring for worker health
- Manual import (create import task without download job)
