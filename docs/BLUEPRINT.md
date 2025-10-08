# Snaggle Blueprint (v2.1)

## Overview
Snaggle is a self-hosted media management platform designed to unify the best features of Sonarr, Radarr, and Overseerr — with a simpler, more transparent design.  
The system emphasizes **user control**, **job transparency**, and **extensibility** over full automation.

---

## Core Philosophy
- Filesystem is the source of truth.
- Postgres maintains metadata, jobs, and user settings.
- Acquisition is flexible: either fully automatic (profiles) or manual (“requires attention” queue).
- Everything runs inside one container: Nginx + Go backend + Postgres + Redis (for caching).

---

## System Architecture

### Components
| Component | Description |
|------------|--------------|
| **Backend (Go + Echo)** | Core API, job queue, acquisition logic. |
| **Frontend (Vue 3 SPA)** | User interface served by Nginx. |
| **Database (Postgres)** | Persistent data (users, libraries, media, jobs). |
| **Cache (Redis)** | Temporary storage for discovery candidates and ephemeral data. |
| **Worker (in-proc)** | Handles background jobs: scanning, discovery, download, import, refresh. |

---

## Key Concepts

### Library
A root folder on disk that Snaggle monitors. Each library has:
- Type: `movie` or `series`
- Root path (validated)
- Enabled flag
- Last scan date

### Media Item
Represents a single movie or series tracked by Snaggle.

### Monitored Item
A “wanted” media entity (movie or episode) that triggers acquisition.

### Discovery
The process of finding possible releases for a monitored item from one or more sources.

### Candidate
A potential release (torrent/NZB). **Ephemeral** — stored in Redis, not persisted long-term.

### Selection
Represents a chosen candidate (manual or automatic).

### Download Task
Tracks a job submitted to a download client.

### Import Job
Moves a completed download into the library and updates the database using a **hardlink-first** strategy (details below).

### Refresh Job
Triggers Plex/Jellyfin refresh after import.

---

## Acquisition Flow

### 1. Discovery
- Triggered manually or periodically.
- Queries provider APIs (Torznab, Prowlarr, etc.).
- Normalizes and scores results.
- Filters out rejected/banned releases.
- **Auto mode:** immediately picks best candidate and proceeds to download.
- **Manual mode:** caches top candidates in Redis for interactive selection.

### 2. Selection
- **Manual:** user chooses from cached list (`POST /api/v1/selections`).  
- **Auto:** chosen automatically based on quality profile rules.

### 3. Download
- Selected release sent to download client (qBittorrent, SABnzbd, etc.).
- Progress tracked by polling or webhook.
- On completion → triggers Import job.

### 4. Import (Hardlink-first strategy)
- **Goal:** keep seeding intact when desired, avoid copies, and be safe across filesystems.
- **Algorithm (pseudocode):**
  1. Ensure source file is finalized by the client (not a temp/incomplete name).  
  2. If source and destination are on the **same filesystem**:
     - Try **hardlink** `src → dst` (preserves inode; seeding continues).  
     - If hardlink fails (permissions/policy), try **rename** `src → dst` (atomic; same inode; only if not seeding).  
  3. If cross-filesystem (EXDEV) or previous steps fail:
     - Try **reflink** (CoW clone) if supported.  
     - Fallback to **copy** with checksum verification.  
  4. Update DB (`media_file`) and fire **Refresh** job.
- **Implications:**
  - Hardlinks require same filesystem; chmod/chown affect all links (same inode).  
  - Rename preserves inode but breaks seeding unless the client is aware.  
  - Reflink gives instant clones but uses a new inode (requires CoW FS).  
  - Copy is universal but slow and uses space.

### 5. Refresh
- If Plex/Jellyfin integrations are configured, refresh library.

---

## Filesystem Layout (recommended)
- Host: `/srv/media` (single filesystem).  
- Container: bind to `/data`.  
- Paths: `/data/downloads/complete/...` and `/data/library/...` on **the same mount** to enable hardlinks/renames.

---

## Database Schema (Core)

```sql
create table app_user (
  id uuid primary key,
  email text unique not null,
  password_hash text,
  role text not null check (role in ('admin','user')),
  created_at timestamptz not null default now()
);

create table library (
  id uuid primary key,
  name text not null,
  type text not null check (type in ('movie','series')),
  root_path text not null,
  enabled boolean not null default true,
  created_at timestamptz not null default now()
);

create table media_item (
  id uuid primary key,
  library_id uuid not null references library(id) on delete cascade,
  type text not null,
  title text not null,
  year int,
  tmdb_id int,
  created_at timestamptz not null default now()
);

create table media_file (
  id uuid primary key,
  media_id uuid not null references media_item(id) on delete cascade,
  path text not null unique,
  size_bytes bigint,
  resolution text,
  added_at timestamptz not null default now()
);

create table monitored_item (
  id uuid primary key,
  media_id uuid not null references media_item(id) on delete cascade,
  type text not null check (type in ('movie','series','episode')),
  desired_profile_id uuid,
  is_active boolean not null default true,
  created_at timestamptz not null default now()
);

create table selection (
  id uuid primary key,
  monitored_id uuid not null references monitored_item(id) on delete cascade,
  candidate_title text not null,
  candidate_link text not null,
  mode text not null check (mode in ('manual','auto')),
  decided_by uuid references app_user(id),
  decided_at timestamptz not null default now()
);

create table download_task (
  id uuid primary key,
  selection_id uuid not null references selection(id) on delete cascade,
  client text not null,
  client_task_id text,
  state text not null check (state in ('queued','downloading','completed','failed')),
  progress numeric(5,2) default 0,
  error text,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);
```

---

## APIs

### Libraries
- `POST /api/v1/libraries`
- `GET /api/v1/libraries`
- `DELETE /api/v1/libraries/:id`
- `POST /api/v1/libraries/:id/scan`

### Media
- `GET /api/v1/library`
- `GET /api/v1/media/:id`

### Monitored Items
- `POST /api/v1/monitored`
- `GET /api/v1/monitored`
- `DELETE /api/v1/monitored/:id`

### Discovery & Selection
- `POST /api/v1/monitored/:id/discover`
- `GET /api/v1/discoveries/:id/candidates`
- `POST /api/v1/selections`
- `POST /api/v1/monitored/:id/ignore`

### Downloads & Jobs
- `GET /api/v1/downloads`
- `GET /api/v1/jobs`
- `GET /api/v1/events` (SSE stream)

---

## Import Helper (Go sketch)
```go
// tryImport does: hardlink → rename → reflink → copy
func tryImport(src, dst string) error {
    if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil { return err }

    sameFS := func(a, b string) bool {
        var sa, sb syscall.Stat_t
        if err := syscall.Stat(a, &sa); err != nil { return false }
        if err := syscall.Stat(filepath.Dir(b), &sb); err != nil { return false }
        return sa.Dev == sb.Dev
    }

    if sameFS(src, dst) {
        if err := os.Link(src, dst); err == nil { return nil }           // hardlink
        if err := os.Rename(src, dst); err == nil { return nil }         // rename (breaks seeding)
    }

    // TODO: try reflink (ioctl FICLONE) when supported; else copy
    return copyFile(src, dst)
}
```

---

## Ephemeral Candidate Storage
- Candidates stored in Redis for ~1 hour (configurable TTL).
- Key: `candidates:{monitored_id}` → array of candidates.
- Value structure:
  ```json
  {
    "title": "Dune.2021.2160p.WEB-DL",
    "quality": "2160p",
    "sizeBytes": 3200000000,
    "seeders": 123,
    "indexer": "Torznab@Jackett",
    "score": 87,
    "link": "magnet:?xt=urn:btih:..."
  }
  ```
- Only **selections** and downstream artifacts are persisted.

---

## Jobs Overview
| Job | Description |
|------|--------------|
| **scan** | Walks library, parses filenames, upserts media records. |
| **discover** | Searches providers, normalizes results, caches candidates. |
| **download** | Sends torrent/NZB to client and monitors progress. |
| **import** | **Hardlink-first** move into library, upserts `media_file`. |
| **refresh** | Triggers Plex/Jellyfin refresh if configured. |

---

## Milestones

### Milestone 1 — Core Foundations
- User auth (JWT, admin seed).
- Libraries CRUD.
- Filesystem scan job.

### Milestone 2 — Monitored Items & Discovery
- Add monitored table + discovery logic.
- Integrate with 1 provider (Torznab).

### Milestone 3 — Manual Selection
- Redis caching for candidates.
- Selection API + download task creation.

### Milestone 4 — Import & Refresh (Hardlink-first)
- Implement hardlink-first importer with safe fallbacks.
- Import job updates media_file and triggers Plex/Jellyfin refresh.

### Milestone 5 — UI / UX
- Vue UI for browsing, requests, queue, requires-attention.

---

## Future Extensions
- Quality profiles (auto mode).
- Multi-provider aggregation.
- Episode-level tracking.
- Advanced file rename templates.
- Notifications (Discord/webhook).
- Mobile-friendly dashboard.
