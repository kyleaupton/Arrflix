#!/usr/bin/env bash
set -euo pipefail

echo "[init-postgres] Creating prowlarr role and databases..."
echo "[init-postgres] POSTGRES_USER: ${POSTGRES_USER:-postgres}"
echo "[init-postgres] POSTGRES_DB: ${POSTGRES_DB:-postgres}"

# Create prowlarr role
psql -v ON_ERROR_STOP=1 --username "${POSTGRES_USER:-postgres}" --no-password <<-'EOSQL'
DO $$ 
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'prowlarr') THEN
    CREATE ROLE prowlarr LOGIN PASSWORD 'prowlarrpw';
  END IF;
END
$$;
EOSQL

# Create prowlarr databases
psql -U "${POSTGRES_USER:-postgres}" -c "CREATE DATABASE \"prowlarr-main\" OWNER \"prowlarr\""
psql -U "${POSTGRES_USER:-postgres}" -c "CREATE DATABASE \"prowlarr-log\" OWNER \"prowlarr\""

echo "[init-postgres] Initialization complete"