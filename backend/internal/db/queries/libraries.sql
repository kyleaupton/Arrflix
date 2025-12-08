-- name: ListLibraries :many
select * from library
order by name asc;

-- name: GetLibrary :one
select * from library
where id = $1;

-- name: CreateLibrary :one
insert into library (name, type, root_path, enabled, "default")
values (sqlc.arg(name), sqlc.arg(type), sqlc.arg(root_path), sqlc.arg(enabled), sqlc.arg(is_default))
returning *;

-- name: UpdateLibrary :one
update library
set name = sqlc.arg(name),
    type = sqlc.arg(type),
    root_path = sqlc.arg(root_path),
    enabled = sqlc.arg(enabled),
    "default" = sqlc.arg(is_default),
    updated_at = now()
where id = sqlc.arg(id)
returning *;

-- name: DeleteLibrary :exec
delete from library where id = $1;

-- name: GetDefaultLibrary :one
select * from library
where type = $1 and "default" = true;


