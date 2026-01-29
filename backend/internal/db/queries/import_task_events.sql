-- Import task events (audit log)

-- name: CreateImportTaskEvent :one
INSERT INTO import_task_event (
  import_task_id,
  event_type,
  old_status,
  new_status,
  message,
  metadata
)
VALUES (
  sqlc.arg(import_task_id),
  sqlc.arg(event_type),
  sqlc.arg(old_status),
  sqlc.arg(new_status),
  sqlc.arg(message),
  sqlc.arg(metadata)
)
RETURNING *;

-- name: ListImportTaskEvents :many
SELECT * FROM import_task_event
WHERE import_task_id = $1
ORDER BY created_at ASC;

-- name: ListRecentImportTaskEvents :many
SELECT * FROM import_task_event
ORDER BY created_at DESC
LIMIT sqlc.arg(limit_val)
OFFSET sqlc.arg(offset_val);

-- name: GetImportTaskTimeline :many
-- Get events for a specific import task
SELECT
  id,
  import_task_id,
  event_type,
  old_status,
  new_status,
  message,
  metadata,
  created_at
FROM import_task_event
WHERE import_task_id = $1
ORDER BY created_at ASC;
