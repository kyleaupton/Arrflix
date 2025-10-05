# North star & constraints

- **North star:** One UI to see what’s in your library, request what’s missing, and watch the pipeline do its thing—without juggling Sonarr/Radarr/Overseerr.
- **Constraints:** Self-hosted, single container, Postgres inside, Echo v4 API, Vue SPA, job runner in-proc for V1, Plex/Jellyfin are **notify-only** (trigger refresh).

# Personas

- **Owner/Admin** (you): manages settings, libraries, approvals, scans.
- **Member**: searches titles, creates requests, sees status.

# V1 outcomes (what “done” means)

- I can add 1+ library roots and scan them.
- I can browse my library with posters/titles/quality at a glance.
- I can submit a request (movie or series), approve it, and see the job run to “imported”.
- I can see a single queue of all jobs with live updates.
- Role-based access: admin vs user.
- Post-import: fire Plex/Jellyfin “refresh” if configured.

# V1 user stories (w/ acceptance criteria)

## Setup & auth

1. **As an admin, I can sign in**

   - Email/password auth; seeded admin on first boot via env.
   - JWT issued; protected routes require it.

2. **As an admin, I can configure libraries**

   - Create library with name, type (`movie|series`), root path.
   - Path exists check; enable/disable flag.

3. **As an admin, I can manage settings**
   - Update global settings: scan interval, watcher on/off, Plex/Jellyfin endpoints.

## Library & scanning

4. **As an admin, I can scan a library**

   - Trigger scan; see job in queue; items appear in Library.
   - Basic parser (Title (Year).ext or SxxEyy.\* for episodes).
   - For series, create seasons/episodes as discovered.

5. **As any user, I can browse/search the library**
   - Grid/list with filters (type, text query) + pagination.
   - Item detail shows files and basic quality (e.g., 1080p).

## Requests

6. **As a user, I can request a title**

   - Search TMDB/TVDB; if missing locally, create request.
   - Shows “Pending approval” unless auto-approve is on.

7. **As an admin, I can approve/deny a request**
   - Approve → enqueue fetch job (demo adapter ok), then import, then appears in Library.
   - Deny → request closed; audit logged.

## Queue & visibility

8. **As any user, I can see pipeline status**
   - Queue view with states: queued → fetching → processing → imported/failed.
   - SSE/WebSocket live updates.

## Integrations (notify-only)

9. **As an admin, I can trigger Plex/Jellyfin refresh after import**
   - On imported event, call configured refresh endpoint (best-effort).

# Backend plan (entities → endpoints → jobs)

## Core entities

- `app_user (id, email, password_hash, role)` – V1 roles: `admin|user` (granular grants later).
- `library (id, name, type, root_path, enabled)`
- `media_item (id, library_id, type, title, year, tmdb_id?, tvdb_id?, created_at)`
- `media_season (id, media_id, season_number)`
- `media_episode (id, season_id, episode_number, title?, air_date?)`
- `media_file (id, media_id, season_id?, episode_id?, path, size_bytes, resolution?, added_at)`
- `request (id, user_id, type, tmdb_id?, tvdb_id?, title, year?, status[pending|approved|denied|fulfilled], created_at)`
- `job_run (id, type[scan|fetch|import|refresh], state[queued|running|done|failed], payload jsonb, started_at, finished_at, error?)`
- `app_setting (key, value jsonb, updated_at)`

## Key endpoints (Echo groups)

- **Auth:**  
  `POST /api/v1/auth/login` → {token}  
  `GET /api/v1/auth/me`
- **Settings (admin):**  
  `GET/PUT /api/v1/settings`
- **Libraries (admin):**  
  `GET/POST /api/v1/libraries`  
  `POST /api/v1/libraries/:id/scan`
- **Library browse:**  
  `GET /api/v1/library?type=&q=&page=&pageSize=`  
  `GET /api/v1/media/:id`
- **Requests:**  
  `POST /api/v1/requests` (user)  
  `GET /api/v1/requests` (user sees theirs; admin sees all)  
  `POST /api/v1/requests/:id/approve` (admin)  
  `POST /api/v1/requests/:id/deny` (admin)
- **Queue/events:**  
  `GET /api/v1/jobs` (recent jobs)  
  `GET /api/v1/events` (SSE)

## Background jobs (MVP)

- **scan**(library*id): walk root, parse files, upsert `media*\*`; emit events.
- **fetch**(request_id): demo/downloader stub; simulate acquisition; on success → enqueue **import**.
- **import**(media_id/path): move/rename (simulated), create `media_file`. On success → enqueue **refresh** if integrations configured.
- **refresh**(service): call Plex/Jellyfin refresh endpoint, log result.

## Scanning pipeline (first pass)

1. Walk filesystem (`.mkv,.mp4` etc.), skip samples.
2. Determine type: movie vs episode (regex on filename/parent folders).
3. Extract `(title, year)` or `(series, season, episode)` from path.
4. **Upsert**:
   - movie: `media_item` (title/year) + `media_file`
   - episode: ensure `media_item` for series, `media_season`, `media_episode`, then `media_file`
5. Resolution heuristic from filename (`2160p/1080p`) → `media_file.resolution`
6. Batch commit per directory to reduce churn.

# Milestones & order of work (sprintable)

### Milestone 0: Foundations (½–1 day)

- Config loader, logger (zerolog), migrations run on boot.
- Health endpoint.

### Milestone 1: Auth & roles (1 day)

- Users table, seed admin from env, login (`/auth/login`), JWT middleware.
- Role gate: admin vs user (single role for V1).

### Milestone 2: Settings & libraries (1 day)

- Settings CRUD (admin), library CRUD with path validation.

### Milestone 3: Scanner (2 days)

- Scan job + endpoint; basic filename parser; upsert `media_*`; library browse API.

### Milestone 4: Queue & SSE (1 day)

- In-proc queue + `/events` stream; `/jobs` history.

### Milestone 5: Requests (2 days)

- TMDB search (basic) + create request; approve/deny; fetch→import stub; library updates; refresh notify.

### Milestone 6: Hardening & DX (1 day)

- Pagination, basic filters, indexes; error model; small e2e happy-path test.

# Acceptance test checklist (V1)

- [ ] Fresh boot seeds admin; can login; protected routes 401 without token.
- [ ] Add a movie library; run scan; files appear in `GET /library`.
- [ ] Create user; user can browse, submit request.
- [ ] Admin approves → job states progress → item imported → shows in library.
- [ ] SSE pushes state changes to the UI in real time.
- [ ] If Plex/Jellyfin configured, “imported” triggers refresh call (200/202 logged).

# Data & index notes (add early)

- `media_item (library_id, type, title, year)` unique index.
- `media_file(path)` unique index; index `(media_id)` for joins.
- `request(status, created_at)` index for dashboards.
- `job_run(created_at)` index for recent queries.

# Risks & mitigations

- **Filename parsing edge cases** → start simple, log “unparsed” files; add overrides later.
- **Long scans** → paginate work; schedule; FS watcher later.
- **Permissions creep** → keep V1 roles simple (admin/user); add granular grants after V1.
