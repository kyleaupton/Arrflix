#!/usr/bin/env bash
set -euo pipefail

: "${DB_HOST:=localhost}"
: "${DB_PORT:=5432}"
: "${DB_USER:=arrflix}"
: "${DB_NAME:=arrflix}"
: "${DB_PASSWORD:=arrflixpw}"

echo "[postgres-ready] checking SQL readiness on ${DB_HOST}:${DB_PORT}/${DB_NAME}..."

for i in {1..20}; do
  if PGPASSWORD="${DB_PASSWORD:-}" psql \
      --host="$DB_HOST" \
      --port="$DB_PORT" \
      --username="$DB_USER" \
      --dbname="$DB_NAME" \
      --command="select 1;" >/dev/null 2>&1; then
    echo "[postgres-ready] DB is ready"
    exit 0
  fi

  echo "[postgres-ready] not ready yet, retrying..."
  sleep 2
done

echo "[postgres-ready] ERROR: postgres never became ready" >&2
exit 1
