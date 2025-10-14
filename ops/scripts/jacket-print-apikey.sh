#!/usr/bin/env bash
set -euo pipefail
conf="${JACKETT_CONFIG:-/root/.config/Jackett/ServerConfig.json}"
for i in {1..60}; do
if [ -s "$conf" ]; then
    jq -r '.APIKey' "$conf"
    exit 0
fi
sleep 1
done
echo "Timeout waiting for $conf" >&2
exit 1
