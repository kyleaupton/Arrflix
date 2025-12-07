#!/usr/bin/env bash
set -euo pipefail

# This script retrieves the Prowlarr API key from the shared environment file
# It's used by the API service to get the API key that was generated for Prowlarr

ENV_FILE="/tmp/prowlarr.env"

# Wait for the environment file to be created by Prowlarr startup script
for i in {1..60}; do
    if [ -f "$ENV_FILE" ] && [ -s "$ENV_FILE" ]; then
        # Source the environment file and extract the API key
        source "$ENV_FILE"
        if [ -n "${PROWLARR_API_KEY:-}" ]; then
            echo "$PROWLARR_API_KEY"
            exit 0
        fi
    fi
    sleep 1
done

echo "Timeout waiting for Prowlarr API key environment file" >&2
exit 1
