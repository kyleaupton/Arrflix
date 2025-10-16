#!/usr/bin/env bash
set -euo pipefail
conf="${PROWLARR_CONFIG:-/root/.config/Prowlarr/config.xml}"
for i in {1..60}; do
if [ -s "$conf" ]; then
    # Prowlarr stores API key in config.xml, extract it using sed/xml parsing
    grep -oP '(?<=<ApiKey>)[^<]+' "$conf" || echo ""
    exit 0
fi
sleep 1
done
echo "Timeout waiting for $conf" >&2
exit 1
