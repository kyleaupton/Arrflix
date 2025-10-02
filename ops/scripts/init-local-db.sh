#!/usr/bin/env bash
set -euo pipefail

# stop and remove any existing containers
docker compose -f docker-compose.db.yml down -v --remove-orphans

# initialize local database
docker compose -f docker-compose.db.yml up -d
