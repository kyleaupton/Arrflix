-- name: ListNameTemplates :many
select * from name_template
order by name asc;

-- name: GetNameTemplate :one
select * from name_template
where id = $1;

-- name: GetDefaultNameTemplate :one
select * from name_template
where type = $1 and "default" = true;

-- name: CreateNameTemplate :one
insert into name_template (name, type, template, "default")
values (sqlc.arg(name), sqlc.arg(type), sqlc.arg(template), sqlc.arg(is_default))
returning *;

-- name: UpdateNameTemplate :one
update name_template
set name = sqlc.arg(name),
    type = sqlc.arg(type),
    template = sqlc.arg(template),
    "default" = sqlc.arg(is_default),
    updated_at = now()
where id = sqlc.arg(id)
returning *;

-- name: DeleteNameTemplate :exec
delete from name_template where id = $1;

