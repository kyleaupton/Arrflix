-- Download jobs

-- name: CreateDownloadJob :one
insert into download_job (
  status,
  protocol,
  media_type,
  media_item_id,
  season_id,
  episode_id,
  indexer_id,
  guid,
  candidate_title,
  candidate_link,
  downloader_id,
  library_id,
  name_template_id,
  predicted_dest_path
)
values (
  'created',
  sqlc.arg(protocol),
  sqlc.arg(media_type),
  sqlc.arg(media_item_id),
  sqlc.arg(season_id),
  sqlc.arg(episode_id),
  sqlc.arg(indexer_id),
  sqlc.arg(guid),
  sqlc.arg(candidate_title),
  sqlc.arg(candidate_link),
  sqlc.arg(downloader_id),
  sqlc.arg(library_id),
  sqlc.arg(name_template_id),
  sqlc.arg(predicted_dest_path)
)
on conflict (indexer_id, guid) do update
set updated_at = now()
returning *;

-- name: GetDownloadJob :one
select * from download_job
where id = $1;

-- name: GetDownloadJobByCandidate :one
select * from download_job
where indexer_id = $1 and guid = $2;

-- name: ListDownloadJobsByMediaItem :many
select * from download_job
where media_item_id = $1
order by created_at desc;

-- name: ListDownloadJobs :many
select * from download_job
order by created_at desc;

-- name: ListDownloadJobsByTmdbMovieID :many
select j.*
from download_job j
join media_item mi on mi.id = j.media_item_id
where mi.type = 'movie' and mi.tmdb_id = $1
order by j.created_at desc;

-- name: ListDownloadJobsByTmdbSeriesID :many
select j.*
from download_job j
join media_item mi on mi.id = j.media_item_id
where mi.type = 'series' and mi.tmdb_id = $1
order by j.created_at desc;

-- name: CancelDownloadJob :one
update download_job
set status = 'cancelled',
    updated_at = now()
where id = $1
returning *;

-- name: SetDownloadJobEnqueued :one
update download_job
set status = 'enqueued',
    downloader_external_id = sqlc.arg(downloader_external_id),
    updated_at = now()
where id = sqlc.arg(id)
returning *;

-- name: SetDownloadJobDownloadSnapshot :one
update download_job
set status = sqlc.arg(status),
    downloader_status = sqlc.arg(downloader_status),
    progress = sqlc.arg(progress),
    download_save_path = sqlc.arg(download_save_path),
    download_content_path = sqlc.arg(download_content_path),
    updated_at = now()
where id = sqlc.arg(id)
returning *;

-- name: SetDownloadJobImporting :one
update download_job
set status = 'importing',
    import_source_path = sqlc.arg(import_source_path),
    updated_at = now()
where id = sqlc.arg(id)
returning *;

-- name: SetDownloadJobImported :one
update download_job
set status = 'imported',
    import_source_path = sqlc.arg(import_source_path),
    import_dest_path = sqlc.arg(import_dest_path),
    import_method = sqlc.arg(import_method),
    primary_media_file_id = sqlc.arg(primary_media_file_id),
    updated_at = now()
where id = sqlc.arg(id)
returning *;

-- name: BumpDownloadJobRetry :one
update download_job
set attempt_count = attempt_count + 1,
    last_error = sqlc.arg(last_error),
    next_run_at = sqlc.arg(next_run_at),
    updated_at = now()
where id = sqlc.arg(id)
returning *;

-- name: MarkDownloadJobFailed :one
update download_job
set status = 'failed',
    last_error = sqlc.arg(last_error),
    updated_at = now()
where id = sqlc.arg(id)
returning *;

-- name: ClaimRunnableDownloadJobs :many
with cte as (
  select id
  from download_job
  where status in ('created','enqueued','downloading','completed','importing')
    and next_run_at <= now()
  order by next_run_at asc
  for update skip locked
  limit $1
)
update download_job j
set updated_at = now()
from cte
where j.id = cte.id
returning j.*;

-- name: LinkDownloadJobMediaFile :exec
insert into download_job_media_file (download_job_id, media_file_id)
values (sqlc.arg(download_job_id), sqlc.arg(media_file_id))
on conflict (download_job_id, media_file_id) do nothing;


