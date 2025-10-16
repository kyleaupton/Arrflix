#!/usr/bin/env bash
set -euo pipefail

: "${PROWLARR_URL:=http://127.0.0.1:9696/}"

echo "[wait-for-prowlarr] Waiting for Prowlarr at ${PROWLARR_URL} ..."
for i in {1..60}; do
  if curl -fsS -o /dev/null "${PROWLARR_URL}"; then
    echo "[wait-for-prowlarr] Prowlarr is ready."
    exit 0
  fi
  echo "[wait-for-prowlarr] Prowlarr is not ready."
  sleep 1
done

echo "[wait-for-prowlarr] Timeout waiting for Prowlarr."
exit 1
