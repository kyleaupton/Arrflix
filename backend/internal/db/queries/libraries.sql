-- name: ListLibraries :many
select * from library
order by name asc;

-- name: GetLibrary :one
select * from library
where id = $1;

-- name: CreateLibrary :one
insert into library (name, type, root_path, enabled)
values (sqlc.arg(name), sqlc.arg(type), sqlc.arg(root_path), sqlc.arg(enabled))
returning *;

-- name: UpdateLibrary :one
update library
set name = $2,
    type = $3,
    root_path = $4,
    enabled = $5,
    updated_at = now()
where id = $1
returning *;

-- name: DeleteLibrary :exec
delete from library where id = $1;


