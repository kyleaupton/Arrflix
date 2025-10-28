# Snaggle

A self-hosted media management platform that unifies the best parts of Sonarr, Radarr, and Overseerr — but with more transparency, flexibility, and developer-first clarity.

## Quick Start

Snaggle now uses a container orchestration approach. Simply run:

```bash
cd ops
docker-compose up -d
```

This starts the `snaggle-ops` reconciler which automatically spawns and manages all required services:

- PostgreSQL database
- Prowlarr indexer management
- Snaggle API backend
- Nginx web server
- Dynamic services (qBittorrent instances, etc.)

Access the web interface at http://localhost:8484

### Development Mode

For development with live reload:

```bash
cd ops
RUNTIME_MODE=dev docker-compose up -d
```

This enables:

- Go API with Air live reload
- Vue frontend with Vite HMR
- Direct access to dev servers on ports 8080 and 5173

## Architecture

Snaggle uses a Go-based reconciliation controller that manages Docker containers in a Kubernetes-style pattern:

- **Single Entry Point**: Only `snaggle-ops` container in docker-compose
- **Dynamic Services**: Add/remove services via database without restarts
- **Better Isolation**: Each service runs in its own container
- **Dependency Management**: Services start in correct order with health checks

See [ops/README.md](ops/README.md) for detailed architecture documentation.

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
- **Resource-aware** (don't clobber the NAS or saturate WAN)

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
