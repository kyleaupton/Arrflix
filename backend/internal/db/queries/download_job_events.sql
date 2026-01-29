-- Download job events (audit log)

-- name: CreateDownloadJobEvent :one
INSERT INTO download_job_event (
  download_job_id,
  event_type,
  old_status,
  new_status,
  message,
  metadata
)
VALUES (
  sqlc.arg(download_job_id),
  sqlc.arg(event_type),
  sqlc.arg(old_status),
  sqlc.arg(new_status),
  sqlc.arg(message),
  sqlc.arg(metadata)
)
RETURNING *;

-- name: ListDownloadJobEvents :many
SELECT * FROM download_job_event
WHERE download_job_id = $1
ORDER BY created_at ASC;

-- name: ListRecentDownloadJobEvents :many
SELECT * FROM download_job_event
ORDER BY created_at DESC
LIMIT sqlc.arg(limit_val)
OFFSET sqlc.arg(offset_val);
