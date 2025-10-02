-- name: ListUserRoles :many
SELECT r.* FROM role r
JOIN user_role ur ON ur.role_id = r.id
WHERE ur.user_id = $1;

-- name: AssignRole :exec
INSERT INTO user_role (user_id, role_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;