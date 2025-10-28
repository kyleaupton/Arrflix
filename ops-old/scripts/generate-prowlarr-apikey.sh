#!/usr/bin/env bash
set -euo pipefail

# Generate a secure 32-character API key using openssl
# This creates a random hex string that's suitable for API authentication
PROWLARR_API_KEY=$(openssl rand -hex 16)

# Export the API key so it's available to child processes
export PROWLARR_API_KEY

# Also export it with the Prowlarr environment variable format
export PROWLARR__AUTH__APIKEY="$PROWLARR_API_KEY"

echo "[generate-prowlarr-apikey] Generated API key: $PROWLARR_API_KEY"

# Execute the command passed as arguments with the API key available
exec "$@"
