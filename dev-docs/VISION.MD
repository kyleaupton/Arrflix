# Vision

## Problem

Existing media managers (e.g., Sonarr/Radarr) feel slow/heavy, hard to customize, and brittle to operate.

## Target Users

- Self-hosters who want predictable, low-latency automation
- Power users who want _composable_ pipelines and clear observability

## Product Principles

- **Local-first** (works offline with queued actions)
- **Observable by default** (events, logs, metrics, traces)
- **Composable** (plugin adapters for indexers, downloaders, metadata)
- **Idempotent** (safe retries; no dupes)
- **Resource-aware** (don’t clobber the NAS or saturate WAN)

## Non-Goals (v1)

- Nothing more than solving my workflow needs
- No rich mobile clients (basic mobile web only)
