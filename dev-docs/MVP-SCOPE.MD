# MVP Scope (v1)

## Personas

- Admin
- User

## Core User Stories

1. As an admin, I can add a series/movie and specify quality profile.
2. The system monitors indexers and enqueues matching releases.
3. A download client integration fetches releases.
4. A post-processor imports files, normalizes names, and updates the library.
5. I can see a live pipeline view and retry/ban items.

## Integrations (v1)

- Indexers: PirateBay, IPTorrents
- Downloaders: qBitorrent
- Metadata (Built-in): TMDB
- Notifiers (Stretch goal): Email

## Quality Profiles (v1)

- Simple: 1080p, 4K, "Best available"
- Single keep policy (no upgrade ladder yet)

## UI (v1)

- Screen to see items in collection
- Screen to add items to collection
- Settings
  - Integrations
  - Naming format
