#!/usr/bin/env bash
set -euo pipefail

: "${POSTGRES_DB:=snaggle}"
: "${POSTGRES_USER:=snaggle}"
: "${POSTGRES_PASSWORD:=snagglepw}"
: "${PGDATA:=/var/lib/postgresql/data}"
: "${TEMP_PORT:=5433}"

echo "[run-postgres] Using PGDATA=${PGDATA}"

# Ensure datadir exists/owned
install -d -m 0700 "$PGDATA"
chown -R postgres:postgres "$PGDATA"

if [ ! -s "$PGDATA/PG_VERSION" ]; then
  echo "[run-postgres] First boot: initializing cluster + creating role/db..."
  initdb -D "$PGDATA"

  echo "[run-postgres] Starting temporary postgres on port ${TEMP_PORT}..."
  pg_ctl -D "$PGDATA" -o "-c listen_addresses=127.0.0.1 -c port=${TEMP_PORT}" -w start
  cleanup_temp() {
    echo "[run-postgres] Stopping temporary postgres..."
    pg_ctl -D "$PGDATA" -m fast -w stop || true
  }
  trap cleanup_temp EXIT

  echo "[run-postgres] Creating role if missing..."
  psql -v ON_ERROR_STOP=1 --username postgres -p ${TEMP_PORT} --no-password <<-'EOSQL'
	DO $$ 
	BEGIN
	  IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'snaggle') THEN
	    CREATE ROLE snaggle LOGIN PASSWORD 'snagglepw';
	  END IF;
	END
	$$;
EOSQL

  echo "[run-postgres] Creating database..."
  psql -U postgres -p ${TEMP_PORT} -c "CREATE DATABASE \"${POSTGRES_DB}\" OWNER \"${POSTGRES_USER}\""

  cleanup_temp
  trap - EXIT
else
  echo "[run-postgres] Existing cluster detected (PG_VERSION present); skipping init."
fi

echo "[run-postgres] Starting postgres (final)..."
exec postgres -D "$PGDATA" -c listen_addresses=127.0.0.1
