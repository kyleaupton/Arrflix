# Iteration 001: Media File System Refactor

**Status**: In Progress
**Started**: 2026-01-29
**Goal**: Separate identity from manifestation, add file existence tracking, import history, and unmatched file workflow

---

## Problem Statement

The media file system was built quickly during PoC and has several issues:

1. **UI lies about file existence** - If a file is deleted from the filesystem, the DB still says it exists
2. **`status` column conflates concerns** - Mixes workflow states (downloading, importing) with file states (available, missing)
3. **Redundant `season_id`** - Can be derived from `episode_id`, adds confusion
4. **No import history** - Can't see how a file got into the library or debug failed imports
5. **No handling for unmatched files** - Scanner skips files it can't identify, leaving users unaware

---

## Design Decisions

### 1. Separate File Identity from File State

**Before:**
```
media_file
├── library_id, media_item_id, season_id, episode_id  (identity)
├── path                                               (location)
├── status (available|missing|downloading|...)         (mixed concerns)
└── added_at
```

**After:**
```
media_file                    (identity + location)
├── library_id, media_item_id, episode_id
├── path
└── created_at

media_file_state              (existence tracking)
├── media_file_id (PK)
├── file_exists
├── file_size
└── last_verified_at
```

### 2. Track Import History

New `media_file_import` table tracks every import attempt:
- Links to `media_file` (nullable for failed pre-creation attempts)
- Links to `download_job` (nullable for scans/manual imports)
- Records method: `hardlink`, `copy`, `scan`, `manual_match`
- Tracks success/failure with error messages

### 3. Unmatched Files Workflow

New `unmatched_file` table for files scanner couldn't identify:
- Stores path, file size, discovery time
- Holds suggested matches (JSONB array from TMDB search)
- Tracks resolution (matched to media_file or dismissed)

### 4. Remove Redundant `season_id`

- Removed from `media_file` table
- Removed from `download_job` table
- Derive via join: `episode → season` when needed

### 5. Enforce Media Item Uniqueness

Added `UNIQUE(type, tmdb_id)` constraint on `media_item` to prevent duplicate entries for the same movie/series.

---

## Schema Changes

### New Tables

| Table | Purpose |
|-------|---------|
| `media_file_state` | Track file existence, size, last verification |
| `media_file_import` | Import history (successes and failures) |
| `unmatched_file` | Queue for files needing manual identification |

### Modified Tables

| Table | Change |
|-------|--------|
| `media_item` | Add `UNIQUE(type, tmdb_id)` |
| `media_file` | Remove `season_id`, remove `status` |
| `download_job` | Remove `season_id` |

### Migration Consolidation

Consolidated 13 migration files into 9 logical groupings:

```
0001_extensions_auth.up.sql      - UUID extension, users, roles, permissions
0002_seed_roles_permissions.up.sql - Seed data
0003_settings.up.sql             - app_setting
0004_library.up.sql              - library
0005_media.up.sql                - media_item, season, episode, file, state, import, unmatched
0006_policy_name_template.up.sql - policy, rule, action, name_template
0007_downloader.up.sql           - downloader
0008_download_job.up.sql         - download_job (without season_id)
0009_api_cache.up.sql            - api_cache
```

---

## New API Endpoints

```
GET    /v1/unmatched-files              - List unmatched files (paginated)
GET    /v1/unmatched-files/:id          - Get unmatched file details
POST   /v1/unmatched-files/:id/match    - Manual match to media item
POST   /v1/unmatched-files/:id/dismiss  - Dismiss (ignore) unmatched file
POST   /v1/unmatched-files/:id/refresh  - Re-run TMDB search for suggestions
```

---

## Files Changed

### Database Layer
- `backend/internal/db/migrations/*` - Consolidated migrations
- `backend/internal/db/queries/media.sql` - New queries for state, import, unmatched
- `backend/internal/db/queries/download_jobs.sql` - Remove season_id

### Repository Layer
- `backend/internal/repo/media.go` - Update for new schema
- `backend/internal/repo/download_jobs.go` - Remove season_id
- `backend/internal/repo/media_file_state.go` - New
- `backend/internal/repo/media_file_import.go` - New
- `backend/internal/repo/unmatched_file.go` - New

### Service Layer
- `backend/internal/service/scan.go` - Create unmatched files, file state
- `backend/internal/service/imports.go` - Create import records
- `backend/internal/service/download_candidates.go` - Remove season_id
- `backend/internal/service/unmatched_files.go` - New service
- `backend/internal/service/service.go` - Register new service

### Handler Layer
- `backend/internal/http/handlers/unmatched_files.go` - New handler
- `backend/internal/http/http.go` - Register new routes

---

## Testing Checklist

- [ ] Nuke Docker container, verify migrations run clean
- [ ] Run library scan on known media, verify `media_file` + `media_file_state` created
- [ ] Run library scan on unknown media, verify `unmatched_file` created
- [ ] Queue a download, verify `media_file_import` record on success
- [ ] Test manual match workflow via API
- [ ] Verify frontend still works with updated API responses

---

## Future Work (Out of Scope)

- xattr identity storage (store TMDB ID in file metadata)
- Rich metadata tables (video/audio stream details)
- File verification scheduling (periodic existence checks)
- Download system refactoring (separate iteration)
