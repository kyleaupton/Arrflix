-- name: GetUserByEmail :one
SELECT * FROM app_user WHERE lower(email) = lower($1) AND is_active = true;

-- name: CreateUser :one
INSERT INTO app_user (email, display_name, password_hash, is_active)
VALUES ($1, $2, $3, true)
RETURNING *;

-- name: UpsertIdentity :one
INSERT INTO user_identity (user_id, provider, subject, username, access_token, refresh_token, token_expires_at, raw)
VALUES ($1, $2, $3, $4, $5, $6, $7, COALESCE($8,'{}'::jsonb))
ON CONFLICT (provider, subject) DO UPDATE
SET user_id=EXCLUDED.user_id, username=EXCLUDED.username, access_token=EXCLUDED.access_token,
    refresh_token=EXCLUDED.refresh_token, token_expires_at=EXCLUDED.token_expires_at, raw=EXCLUDED.raw
RETURNING *;

-- name: UpdateUserPassword :exec
UPDATE app_user
SET password_hash = $2, updated_at = now()
WHERE id = $1;