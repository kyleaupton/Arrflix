#!/usr/bin/env bash
set -euo pipefail

# ----- Defaults (can be overridden by container envs) -----
: "${PORT:=8080}"                              # internal API port
: "${POSTGRES_DB:=snaggle}"
: "${POSTGRES_USER:=snaggle}"
: "${POSTGRES_PASSWORD:=snagglepw}"
: "${SSE_ALLOW_ORIGIN:=*}"
# If DATABASE_URL not provided, derive one that talks to local Postgres
: "${DATABASE_URL:=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@127.0.0.1:5432/${POSTGRES_DB}?sslmode=disable}"

export PORT POSTGRES_DB POSTGRES_USER POSTGRES_PASSWORD SSE_ALLOW_ORIGIN DATABASE_URL

echo "[entrypoint] Using DATABASE_URL=${DATABASE_URL}"
echo "[entrypoint] Starting supervisord (Postgres -> wait-for-db -> API -> Nginx)"

exec /usr/bin/supervisord -n -c /etc/supervisor/conf.d/supervisord.conf
