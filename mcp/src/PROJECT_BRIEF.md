# Arrflix Project Brief

## Overview

Arrflix is a self-hosted media management platform (Go/Vue 3) designed to unify the functionality of Sonarr, Radarr, and Overseerr. It prioritizes filesystem transparency, flexible monitoring (Series, Season, or Episode level), and a hardlink-first import strategy.

**Note**: This brief is the primary source of truth for the project's current state. The root `README.md` contains basic user information, but this document should be used for development context.

### Core Philosophy

- **Filesystem is Truth**: The database reflects the state of the disk, not the other way around.
- **Transparency**: Clear visibility into why candidates were selected or ignored.
- **Hardlink-First**: Avoid unnecessary copies; use hardlinks or reflinks whenever possible.

## Getting Started (Development)

The project uses a multi-container Docker setup defined in the root `docker-compose.yml`.

1. **Prerequisites**: Docker and Docker Compose.
2. **Environment**:
   - Set `TMDB_API_KEY` in your shell or an `.env` file.
   - Set `MEDIA_LIBRARIES` to the root path of your media collection.
3. **Run Dev Environment**:
   ```bash
   docker compose up --build
   ```
4. **Access**:
   - **Web UI**: `http://localhost:8484` (Proxied via Nginx)
   - **Prowlarr**: `http://localhost:9697` (Bundled in the main container)
   - **Postgres**: `localhost:5432`
   - **Qbittorrent**: `http://localhost:8485` (Routed through a VPN container)

## Architecture

Arrflix uses a container-based architecture managed by **s6-overlay** for internal process orchestration within the main container.

### Services & Containers

- **Main Container (`arrflix`)**:
  - **Backend (Go + Echo)**: REST API and background worker (in-proc). Uses `sqlc` for type-safe database access and `pgx/v5` for Postgres.
  - **Frontend (Vue 3 + Vite)**: Single Page Application with a generated TypeScript client.
  - **Database (Postgres)**: Persistent storage for metadata, jobs, and settings.
  - **Nginx**: Serves the frontend and proxies API requests.
  - **Prowlarr**: Integrated for indexer management and discovery.
  - **s6-overlay**: Manages the lifecycle of the services above, ensuring correct startup order (e.g., Postgres -> Backend).

### Persistence & Cache

- **Postgres-backed Cache**: Ephemeral data and external API responses are stored in the `api_cache` table in Postgres, replacing the need for Redis.
- **Hardlink Imports**: The system attempts to hardlink files from the download directory to the library directory to save space and reduce IO.

### Key Concepts

- **Library**: Managed root folders for media.
- **Media Item**: A Movie or Series record.
- **Monitored Item**: A subscription to a Series, Season, or Episode for discovery.
- **Service Instance**: Dynamic configuration for external downloaders or indexers.
- **Download Job**: Tracks a release from selection through to final import.
- **Policy & Name Template**: Configurable rules for selecting candidates and naming files.

## Current Status

- [x] **Milestone 1 (Core)**: Auth, settings, and library scanning.
- [x] **Milestone 2 (Series)**: Multi-level monitoring hierarchy.
- [x] **Milestone 3 (Discovery)**: Provider integration and candidate scoring.
- [wip] **Milestone 4 (Import)**: Hardlink-first import worker implementation.
- [ ] **Milestone 5 (UI/UX)**: Advanced dashboard and mobile optimization.

