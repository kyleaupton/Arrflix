# Snaggle

## Vision

### Problem

Existing media managers (e.g., Sonarr/Radarr) feel slow/heavy, hard to customize, and brittle to operate.

### Target Users

- Self-hosters who want predictable, low-latency automation
- Power users who want _composable_ pipelines and clear observability

### Product Principles

- **Local-first** (works offline with queued actions)
- **Observable by default** (events, logs, metrics, traces)
- **Composable** (plugin adapters for indexers, downloaders, metadata)
- **Idempotent** (safe retries; no dupes)
- **Resource-aware** (don’t clobber the NAS or saturate WAN)

### Non-Goals (v1)

- Nothing more than solving my workflow needs
- No rich mobile clients (basic mobile web only)

## MVP Scope (v1)

### Personas

- Admin
- User

### Core User Stories

1. As an admin, I can add a series/movie and specify quality profile.
2. The system monitors indexers and enqueues matching releases.
3. A download client integration fetches releases.
4. A post-processor imports files, normalizes names, and updates the library.
5. I can see a live pipeline view and retry/ban items.

### Integrations (v1)

- Indexers: PirateBay, IPTorrents
- Downloaders: qBitorrent
- Metadata (Built-in): TMDB
- Notifiers (Stretch goal): Email

### Quality Profiles (v1)

- Simple: 1080p, 4K, "Best available"
- Single keep policy (no upgrade ladder yet)

### UI (v1)

- Screen to see items in collection
- Screen to add items to collection
- Settings
  - Integrations
  - Naming format
