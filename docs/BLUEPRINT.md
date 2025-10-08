# Snaggle Blueprint (v2.2)

## Overview

Snaggle is a self-hosted media management platform that unifies the best parts of Sonarr, Radarr, and Overseerr — but with more transparency, flexibility, and developer-first clarity.

The system supports both **movies** and **series**, with **per-episode monitoring**, **interactive acquisitions**, and a **hardlink-first import strategy** to optimize filesystem usage.

---

## Core Philosophy

- Filesystem is the source of truth.
- Database tracks metadata, jobs, and user settings.
- Discovery is user-controllable: automatic or interactive.
- Avoid unnecessary file copies — prefer hardlinks and renames.
- Multi-level monitoring: series, season, or episode.

---

## System Architecture

| Component                | Description                                                |
| ------------------------ | ---------------------------------------------------------- |
| **Backend (Go + Echo)**  | API, job runner, and orchestration logic.                  |
| **Frontend (Vue 3 SPA)** | Web UI for managing media and downloads.                   |
| **Database (Postgres)**  | Persistent store for users, libraries, and jobs.           |
| **Cache (Redis)**        | Ephemeral store for discovery results and transient state. |
| **Worker (in-proc)**     | Handles scanning, discovery, import, and refresh jobs.     |

---

## Hierarchy Overview

```
Library → Media Item (Movie or Series)
Series → Season → Episode
Episode ↔ Monitored Item ↔ Candidate ↔ Selection → Download → Import → File
```

---

## Schema Overview

### Libraries

Represents folders being managed and scanned.

```sql
create table library (
  id uuid primary key,
  name text not null,
  type text not null check (type in ('movie','series')),
  root_path text not null,
  enabled boolean not null default true,
  created_at timestamptz not null default now()
);
```

### Media Items (Movies or Series)

```sql
create table media_item (
  id uuid primary key,
  library_id uuid not null references library(id) on delete cascade,
  type text not null check (type in ('movie','series')),
  title text not null,
  year int,
  tmdb_id int,
  created_at timestamptz not null default now()
);
```

### Seasons and Episodes

```sql
create table media_season (
  id uuid primary key,
  media_id uuid not null references media_item(id) on delete cascade,
  season_number int not null,
  title text,
  air_date date,
  unique (media_id, season_number)
);

create table media_episode (
  id uuid primary key,
  season_id uuid not null references media_season(id) on delete cascade,
  episode_number int not null,
  title text,
  air_date date,
  tmdb_id bigint,
  tvdb_id bigint,
  unique (season_id, episode_number)
);
```

### Media Files

```sql
create table media_file (
  id uuid primary key,
  media_id uuid not null references media_item(id) on delete cascade,
  season_id uuid references media_season(id) on delete set null,
  episode_id uuid references media_episode(id) on delete set null,
  path text not null unique,
  size_bytes bigint,
  resolution text,
  added_at timestamptz not null default now()
);
```

### Monitored Items (multi-level support)

```sql
create table monitored_item (
  id uuid primary key,
  media_id uuid not null references media_item(id) on delete cascade,
  season_id uuid references media_season(id) on delete cascade,
  episode_id uuid references media_episode(id) on delete cascade,
  desired_profile_id uuid,
  is_active boolean not null default true,
  created_at timestamptz not null default now(),
  constraint monitored_target_oneof check (
    (episode_id is not null and season_id is not null and media_id is not null) or
    (episode_id is null and season_id is not null and media_id is not null) or
    (episode_id is null and season_id is null and media_id is not null)
  )
);
```

---

## Import Flow (Hardlink-first)

When importing completed downloads:

1. Verify the client has finalized the file.
2. **If same filesystem:**
   - Attempt **hardlink** (`os.Link(src, dst)`).
   - If hardlink fails, attempt **rename** (safe only if not seeding).
3. **If cross-filesystem:**
   - Attempt **reflink** (copy-on-write).
   - Fallback to **copy** with verification.
4. Update `media_file` record and trigger refresh.

Example Go sketch:

```go
func tryImport(src, dst string) error {
    if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil { return err }

    sameFS := func(a, b string) bool {
        var sa, sb syscall.Stat_t
        if err := syscall.Stat(a, &sa); err != nil { return false }
        if err := syscall.Stat(filepath.Dir(b), &sb); err != nil { return false }
        return sa.Dev == sb.Dev
    }

    if sameFS(src, dst) {
        if err := os.Link(src, dst); err == nil { return nil }
        if err := os.Rename(src, dst); err == nil { return nil }
    }

    return copyFile(src, dst)
}
```

---

## Monitoring Flow

| Level       | Behavior                                |
| ----------- | --------------------------------------- |
| **Series**  | Tracks all future and missing episodes. |
| **Season**  | Tracks only that season's episodes.     |
| **Episode** | Tracks a single episode.                |

Each monitored item spawns discovery jobs scoped by its level:

- Series → all missing episodes.
- Season → missing episodes in that season.
- Episode → exactly one item.

---

## Discovery & Selection

Candidates are ephemeral and stored in Redis for about 1 hour.  
Each discovery job can operate at any monitoring level (series, season, or episode).

- **Auto mode**: pick highest score candidate and proceed.
- **Manual mode**: cache top N results and flag monitored item as `requires_attention`.

Redis key format: `candidates:{monitored_id}`

---

## Job Overview

| Job          | Description                                            |
| ------------ | ------------------------------------------------------ |
| **scan**     | Walk library and update media/seasons/episodes/files.  |
| **discover** | Query providers, normalize results, store in Redis.    |
| **download** | Send selected release to client.                       |
| **import**   | Hardlink-first move into library, upsert `media_file`. |
| **refresh**  | Notify Plex/Jellyfin.                                  |

---

## Milestones

### Milestone 1 — Core System

- Auth & settings store.
- Library CRUD + scanning job.

### Milestone 2 — Series Support

- Add `media_season` and `media_episode`.
- Extend scanner to populate per-episode records.

### Milestone 3 — Monitoring & Discovery

- Multi-level monitoring.
- Candidate caching & discovery logic.

### Milestone 4 — Import & Refresh

- Hardlink-first import system.
- Plex/Jellyfin integration.

### Milestone 5 — UI/UX

- Unified dashboard for media, queue, and attention items.

---

## Future Extensions

- Quality profiles (auto mode).
- Season pack handling.
- Multi-provider merging.
- Advanced file renaming.
- Notifications (Discord/webhook).
- Mobile-friendly dashboard.
