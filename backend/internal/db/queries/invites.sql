-- name: CreateInvite :one
INSERT INTO user_invite (email, invited_by)
VALUES ($1, $2)
RETURNING *;

-- name: GetInviteByEmail :one
SELECT * FROM user_invite
WHERE lower(email) = lower($1) AND claimed_at IS NULL;

-- name: ClaimInvite :exec
UPDATE user_invite
SET claimed_at = now()
WHERE id = $1 AND claimed_at IS NULL;

-- name: ListInvites :many
SELECT * FROM user_invite
ORDER BY created_at DESC;

-- name: DeleteInvite :exec
DELETE FROM user_invite WHERE id = $1;
