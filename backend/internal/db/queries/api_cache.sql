-- name: GetApiCache :one
select key, category, response, status, content_type, headers, stored_at, expires_at
from api_cache
where key = $1
  and expires_at > now();

-- name: UpsertApiCache :exec
insert into api_cache (key, category, response, status, content_type, headers, expires_at)
values (sqlc.arg(key), sqlc.arg(category), sqlc.arg(response), sqlc.arg(status), sqlc.arg(content_type), sqlc.arg(headers), sqlc.arg(expires_at))
on conflict (key)
do update set
	category = excluded.category,
	response = excluded.response,
	status = excluded.status,
	content_type = excluded.content_type,
	headers = excluded.headers,
	stored_at = now(),
	expires_at = excluded.expires_at;

-- name: DeleteExpiredApiCache :exec
delete from api_cache where expires_at <= now();