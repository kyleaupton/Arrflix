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

:
: "${RUNTIME_MODE:=prod}"
echo "[entrypoint] Using DATABASE_URL=${DATABASE_URL}"
echo "[entrypoint] RUNTIME_MODE=${RUNTIME_MODE}"

if [ "$RUNTIME_MODE" = "dev" ]; then
  echo "[entrypoint] Using dev supervisor + nginx configs"
  # Replace default nginx config with dev proxy
  cp -f /etc/nginx/conf.d/default.dev.conf /etc/nginx/conf.d/default.conf || true
  exec /usr/bin/supervisord -n -c /etc/supervisor/conf.d/supervisord.dev.conf
else
  echo "[entrypoint] Using prod supervisor + nginx configs"
  exec /usr/bin/supervisord -n -c /etc/supervisor/conf.d/supervisord.conf
fi
