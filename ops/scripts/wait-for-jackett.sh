#!/usr/bin/env bash
set -euo pipefail

: "${JACKETT_URL:=http://127.0.0.1:9117/}"

echo "[wait-for-jackett] Waiting for Jackett at ${JACKETT_URL} ..."
for i in {1..60}; do
  if curl -fsS -o /dev/null "${JACKETT_URL}"; then
    echo "[wait-for-jackett] Jackett is ready."
    exit 0
  fi
  echo "[wait-for-jackett] Jackett is not ready."
  sleep 1
done

echo "[wait-for-jackett] Timeout waiting for Jackett."
exit 1
