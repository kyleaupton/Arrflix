# Iteration 002: Integration Abstraction Layer

**Status**: Complete
**Started**: 2026-01-29
**Goal**: Isolate third-party integrations (Prowlarr) at the boundary with validation and domain types

---

## Problem Statement

The Prowlarr integration was tightly coupled throughout the codebase:

1. **Type leakage** - `prowlarr.Search` types from gostarr library used directly in services
2. **No validation at boundary** - Invalid results (empty `DownloadURL`) caused errors downstream in the downloader
3. **Scattered mapping logic** - `searchResultToCandidate()` buried in `DownloadCandidatesService`
4. **Not swappable** - Can't replace Prowlarr without touching multiple files
5. **Hard to test** - Unit testing requires Prowlarr instance or complex mocking

---

## Design Decisions

### 1. Define Domain Types

Created Arrflix-owned types that don't depend on external libraries:

```go
// SearchQuery - what we ask for
type SearchQuery struct {
    Query     string
    MediaType MediaType  // "movie" or "series"
    Season    *int
    Episode   *int
    Limit     int
}

// SearchResult - what we get back (validated)
type SearchResult struct {
    IndexerID   int64
    IndexerName string
    GUID        string
    Title       string     // guaranteed non-empty
    DownloadURL string     // guaranteed non-empty
    Protocol    string     // "torrent" or "usenet"
    // ... metadata fields
}
```

### 2. Interface for Swappability

```go
type IndexerSource interface {
    Search(ctx context.Context, query SearchQuery) ([]SearchResult, error)
    ListIndexers(ctx context.Context) ([]IndexerInfo, error)
    Test(ctx context.Context) error
}
```

### 3. Adapter Pattern

`ProwlarrSource` implements `IndexerSource`:
- Maps `SearchQuery` → `prowlarr.SearchInput`
- Executes search via Prowlarr client
- Maps `prowlarr.Search` → `SearchResult` with validation
- Filters out invalid results (logs at debug level)

### 4. Validate at Boundary

The adapter handles Prowlarr quirks and validates before results enter the system:

```go
// Fallback chain for download URL
1. r.DownloadURL          // Direct download (IPTorrents, etc.)
2. r.GUID (if magnet:)    // Magnet with trackers (TPB, etc.)
3. Construct from InfoHash // Bare magnet (last resort)

// Reject if no URL found or title empty
```

### 5. Keep IndexerService for Config

`IndexerService` retains CRUD operations for indexer configuration. Only the `Search()` method moved to the adapter.

---

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│ internal/indexer/                                       │
│                                                         │
│   types.go          ← SearchQuery, SearchResult, etc.  │
│   source.go         ← IndexerSource interface          │
│                                                         │
│   prowlarr/                                            │
│     adapter.go      ← ProwlarrSource implementation    │
└─────────────────────────────────────────────────────────┘

Data flow:
  DownloadCandidatesService
       ↓
  IndexerSource.Search(SearchQuery)
       ↓
  ProwlarrSource (adapter)
       ↓
  prowlarr.Prowlarr.SearchContext()
       ↓
  []*prowlarr.Search (raw results)
       ↓
  mapResult() with validation
       ↓
  []SearchResult (clean domain types)
```

---

## Files Created

| File | Purpose |
|------|---------|
| `backend/internal/indexer/types.go` | Domain types: `SearchQuery`, `SearchResult`, `IndexerInfo`, `MediaType` |
| `backend/internal/indexer/source.go` | `IndexerSource` interface definition |
| `backend/internal/indexer/prowlarr/adapter.go` | `ProwlarrSource` implementation with validation |

---

## Files Modified

| File | Change |
|------|--------|
| `backend/internal/service/download_candidates.go` | Use `IndexerSource` instead of `*IndexerService`, cache `SearchResult` instead of `*prowlarr.Search` |
| `backend/internal/service/indexer.go` | Remove `Search()` method, add `Client()` accessor for adapter |
| `backend/internal/service/service.go` | Create `ProwlarrSource`, inject into `DownloadCandidatesService` |

---

## Key Implementation Details

### URL Fallback Logic

Prowlarr returns download URLs in different ways depending on the indexer:

**Private trackers (IPTorrents, etc.):**
```json
{
  "downloadUrl": "https://prowlarr.example.org/1/download?...",
  "guid": "https://iptorrents.com/t/4046305"
}
```

**Public trackers (TPB, etc.):**
```json
{
  "downloadUrl": "",
  "guid": "magnet:?xt=urn:btih:F48D3B26...&tr=...",
  "infoHash": "F48D3B26BF8ED798722964EFF78171D3AAC5B475"
}
```

The adapter handles both:
1. Use `downloadUrl` if present (preferred)
2. Use `guid` if it's a magnet link (has trackers)
3. Construct magnet from `infoHash` (bare, no trackers)
4. Reject if none available

### Seeders/Leechers as Pointers

For usenet, seeders/leechers are meaningless (always 0). Using `*int` allows distinguishing "not applicable" from "zero seeders":

```go
var seeders, leechers *int
if r.Seeders > 0 || string(r.Protocol) == "torrent" {
    seeders = &r.Seeders
}
```

---

## Testing Checklist

- [x] Backend compiles
- [ ] Search movie via API, verify results returned
- [ ] Search series via API, verify results returned
- [ ] Verify invalid results filtered (check debug logs)
- [ ] Preview candidate, verify policy evaluation works
- [ ] Enqueue download, verify job created with valid link
- [ ] Test with private tracker (uses `downloadUrl`)
- [ ] Test with public tracker (uses `guid` magnet)

---

## Future Work (Out of Scope)

- Alternative indexer sources (DirectTorznab, Jackett)
- Unit tests for `mapResult()` edge cases
- Indexer config UI improvements
- Caching strategy changes (currently 5-minute in-memory)
