-- name: GetSetting :one
SELECT key, type, value_json, version, updated_at FROM app_setting WHERE key = $1;

-- name: ListSettings :many
SELECT key, type, value_json, version, updated_at FROM app_setting ORDER BY key;

-- name: UpsertSetting :exec
INSERT INTO app_setting (key, type, value_json)
VALUES ($1, $2, $3)
ON CONFLICT (key) DO UPDATE SET type = EXCLUDED.type, value_json = EXCLUDED.value_json, version = app_setting.version + 1, updated_at = now();


