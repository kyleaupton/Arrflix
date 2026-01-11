-- name: GetSystemInitialized :one
SELECT (value_json)::bool as initialized
FROM app_setting
WHERE key = 'system.initialized';

-- name: SetSystemInitialized :exec
UPDATE app_setting
SET value_json = 'true'::jsonb,
    updated_at = now(),
    version = version + 1
WHERE key = 'system.initialized'
  AND (value_json)::bool = false;

-- name: CountUsers :one
SELECT COUNT(*) FROM app_user;
