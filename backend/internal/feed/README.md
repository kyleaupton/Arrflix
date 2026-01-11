# Home Feed System – Overview

We are implementing a Netflix-like home feed for movies and TV series
using TMDB as the primary content source.

This system must:

- Support movies AND TV series
- Work without any user personalization
- Improve when user signals exist (local or Plex SSO)
- Avoid content repetition across rows
- Feel curated and dynamic, not like static category lists

This system must NOT:

- Require Plex data to function
- Embed personalization logic directly into TMDB fetch code
- Hardcode rows directly in controllers
- Introduce ML or recommendation engines

The backend is responsible for composing the feed.
The frontend simply renders the feed.

Implementation should prioritize:

- Clear abstractions
- Cache-friendly TMDB access
- Incremental extensibility

# Home Feed Mental Model

The home screen is a composed feed, not a set of independent rails.

A Home Feed consists of:

- One hero item
- An ordered list of rows

Rows are selected, populated, deduplicated, and ranked on the backend.

Each row:

- Has a semantic purpose (e.g. “Trending This Week”)
- Pulls from one or more content sources
- Applies deduplication and diversity rules
- Produces a fixed number of items

Movies and TV series are treated uniformly as "Titles".

# Row System

Rows are defined declaratively.

A row definition includes:

- ID
- Display title and subtitle
- Content kind (movie, tv, mixed)
- Source definition (how candidates are fetched)
- Target size and fetch size
- Ranking strategy
- Optional diversity rules
- Optional dependency on user signals

Row definitions are registered centrally and selected dynamically.

Row source providers:

- Only fetch candidate IDs and basic metadata
- Do NOT handle deduplication
- Do NOT apply personalization

# Deduplication & Diversity Rules

Deduplication is applied at composition time.

A global dedupe set is maintained while building the feed.
Rows may also define their own local dedupe scope.

Diversity rules may include:

- Max items per genre per row
- Max items per language per row
- Max items from same collection/franchise

Deduplication and diversity are enforced while selecting items,
not after the row is already full.

# User Signals & Personalization

The system must function without any user signals.

If available, user signals may include:

- Watch history
- Likes / dislikes
- Watch progress
- Imported Plex events

Signals are normalized into internal events.
No code should depend directly on Plex APIs.

Personalization influences:

- Row selection
- Ranking within rows

Personalization does NOT:

- Change row definitions
- Bypass deduplication or diversity rules

# Home Feed Data Flow

1. Build request context (user, region, signals)
2. Select rows from registry
3. Fetch candidates for each row (cached)
4. Compose rows:
   - apply global dedupe
   - apply diversity constraints
   - rank items
5. Hydrate final items in bulk
6. Return feed payload

# Non-Goals

- No recommendation ML
- No real-time learning loops
- No background jobs required
- No frontend business logic
- No dependency on Plex availability
