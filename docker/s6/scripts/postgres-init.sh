#!/usr/bin/env bash
set -euo pipefail

: "${POSTGRES_DB:=snaggle}"
: "${POSTGRES_USER:=snaggle}"
: "${POSTGRES_PASSWORD:=snagglepw}"
: "${PGDATA:=/var/lib/postgresql/data}"
: "${TEMP_PORT:=5433}"
: "${PG_BINDIR:=/usr/lib/postgresql/16/bin}"

export PATH="$PG_BINDIR:$PATH"

echo "[postgres-init] Using PGDATA=${PGDATA}"

install -d -m 0700 "$PGDATA"
chown -R postgres:postgres "$PGDATA"

if [ ! -s "$PGDATA/PG_VERSION" ]; then
  echo "[postgres-init] First boot: initializing cluster + creating role/db..."
  s6-setuidgid postgres initdb -D "$PGDATA"

  # echo "[postgres-init] Setting pg_hba.conf to trust local connections..."
  # echo "local all all trust" >> "$PGDATA/pg_hba.conf"
  # echo "host all all 127.0.0.1/32 trust" >> "$PGDATA/pg_hba.conf"
  # echo "host all all ::1/128 trust" >> "$PGDATA/pg_hba.conf"

  echo "[postgres-init] Setting pg_hba.conf to trust all connections..."
  echo "host all all 0.0.0.0/0 trust" >> "$PGDATA/pg_hba.conf"
  echo "host all all ::/0 trust" >> "$PGDATA/pg_hba.conf"

  echo "[postgres-init] Starting temporary postgres on port ${TEMP_PORT}..."
  s6-setuidgid postgres pg_ctl -D "$PGDATA" \
    -o "-c listen_addresses=127.0.0.1 -c port=${TEMP_PORT}" \
    -w start

  cleanup_temp() {
    echo "[postgres-init] Stopping temporary postgres..."
    s6-setuidgid postgres pg_ctl -D "$PGDATA" -m fast -w stop || true
  }
  trap cleanup_temp EXIT

  echo "[postgres-init] Creating roles if missing..."
  s6-setuidgid postgres psql -v ON_ERROR_STOP=1 -p "${TEMP_PORT}" <<-'EOSQL'
    DO $$
    BEGIN
      IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'snaggle') THEN
        CREATE ROLE snaggle LOGIN PASSWORD 'snagglepw';
      END IF;
      IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'prowlarr') THEN
        CREATE ROLE prowlarr LOGIN PASSWORD 'prowlarrpw';
      END IF;
    END
    $$;
EOSQL

  echo "[postgres-init] Creating databases..."
  s6-setuidgid postgres psql -p "${TEMP_PORT}" -c "CREATE DATABASE \"${POSTGRES_DB}\" OWNER \"${POSTGRES_USER}\""
  s6-setuidgid postgres psql -p "${TEMP_PORT}" -c "CREATE DATABASE \"prowlarr-main\" OWNER \"prowlarr\""
  s6-setuidgid postgres psql -p "${TEMP_PORT}" -c "CREATE DATABASE \"prowlarr-log\" OWNER \"prowlarr\""

  cleanup_temp
  trap - EXIT
else
  echo "[postgres-init] Existing cluster detected; skipping init."
fi

echo "[postgres-init] Done."
