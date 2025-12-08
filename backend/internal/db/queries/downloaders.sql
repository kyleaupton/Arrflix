-- Downloaders

-- name: ListDownloaders :many
select * from downloader
order by name asc;

-- name: GetDownloader :one
select * from downloader
where id = $1;

-- name: GetDefaultDownloader :one
select * from downloader
where protocol = $1 and "default" = true;

-- name: CreateDownloader :one
insert into downloader (name, type, protocol, url, username, password, config_json, enabled, "default")
values (sqlc.arg(name), sqlc.arg(downloader_type), sqlc.arg(protocol), sqlc.arg(url), sqlc.arg(username), sqlc.arg(password), sqlc.arg(config_json), sqlc.arg(enabled), sqlc.arg(is_default))
returning *;

-- name: UpdateDownloader :one
update downloader
set name = sqlc.arg(name),
    type = sqlc.arg(downloader_type),
    protocol = sqlc.arg(protocol),
    url = sqlc.arg(url),
    username = sqlc.arg(username),
    password = sqlc.arg(password),
    config_json = sqlc.arg(config_json),
    enabled = sqlc.arg(enabled),
    "default" = sqlc.arg(is_default),
    updated_at = now()
where id = sqlc.arg(id)
returning *;

-- name: DeleteDownloader :exec
delete from downloader where id = $1;

