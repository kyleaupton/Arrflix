#!/usr/bin/env bash
set -euo pipefail

# Generate a secure 32-character API key
PROWLARR_API_KEY=$(openssl rand -hex 16)
PROWLARR_PORT="9696"

# Prowlarr database configuration
PROWLARR_DB_HOST="localhost"
PROWLARR_DB_PORT="5432"
PROWLARR_DB_NAME="prowlarr"
PROWLARR_DB_USER="prowlarr"
PROWLARR_DB_PASSWORD="prowlarrpw"

# Create a shared environment file that other processes can read
ENV_FILE="/tmp/prowlarr.env"
echo "PROWLARR_API_KEY=$PROWLARR_API_KEY" > "$ENV_FILE"
echo "PROWLARR__AUTH__APIKEY=$PROWLARR_API_KEY" >> "$ENV_FILE"
echo "PROWLARR__SERVER__PORT=$PROWLARR_PORT" >> "$ENV_FILE"

# Make the file readable by all users (since processes might run as different users)
chmod 644 "$ENV_FILE"

echo "[start-prowlarr-with-apikey] Generated API key: $PROWLARR_API_KEY"
echo "[start-prowlarr-with-apikey] Environment file created at: $ENV_FILE"

# Export the API key for the current process
export PROWLARR_API_KEY
export PROWLARR__AUTH__APIKEY="$PROWLARR_API_KEY"
export PROWLARR__SERVER__PORT="$PROWLARR_PORT"
export PROWLARR__AUTH__METHOD="External"
# export PROWLARR__AUTH__ENABLED="false"
# export PROWLARR__AUTH__REQUIRED="false"

# DB config
export PROWLARR__POSTGRES__HOST="$PROWLARR_DB_HOST"
export PROWLARR__POSTGRES__PORT="$PROWLARR_DB_PORT"
export PROWLARR__POSTGRES__USER="$PROWLARR_DB_USER"
export PROWLARR__POSTGRES__PASSWORD="$PROWLARR_DB_PASSWORD"

# Start Prowlarr with the API key
exec /opt/prowlarr/Prowlarr -nobrowser -data=/var/lib/prowlarr/ "$@"
