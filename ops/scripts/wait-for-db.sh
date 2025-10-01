#!/usr/bin/env bash
set -euo pipefail

# Default to same as official image envs if user didn't specify DATABASE_URL.
: "${POSTGRES_DB:=snaggle}"
: "${POSTGRES_USER:=snaggle}"
: "${POSTGRES_PASSWORD:=snagglepw}"
: "${DATABASE_URL:=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@127.0.0.1:5432/${POSTGRES_DB}?sslmode=disable}"

echo "[wait-for-db] Waiting for Postgres at 127.0.0.1:5432 (${POSTGRES_DB}) ..."
for i in {1..60}; do
  if pg_isready -h 127.0.0.1 -p 5432 -d "${POSTGRES_DB}" >/dev/null 2>&1; then
    echo "[wait-for-db] Postgres is ready."
    exit 0
  fi
  sleep 1
done

echo "[wait-for-db] Timeout waiting for Postgres."
exit 1
