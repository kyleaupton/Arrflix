-- Media core schema: media_item, media_season, media_episode, media_file, media_file_state, media_file_import, unmatched_file

-- Media items (movies or series) with unique constraint on type+tmdb_id
CREATE TABLE IF NOT EXISTS media_item (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  type TEXT NOT NULL CHECK (type IN ('movie','series')),
  title TEXT NOT NULL,
  year INT,
  tmdb_id BIGINT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (type, tmdb_id)
);

CREATE INDEX IF NOT EXISTS idx_media_item_type ON media_item (type);
CREATE INDEX IF NOT EXISTS idx_media_item_tmdb ON media_item (tmdb_id);

-- Seasons (only for series)
CREATE TABLE IF NOT EXISTS media_season (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  media_item_id UUID NOT NULL REFERENCES media_item(id) ON DELETE CASCADE,
  season_number INT NOT NULL,
  air_date DATE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (media_item_id, season_number)
);

CREATE INDEX IF NOT EXISTS idx_media_season_media ON media_season (media_item_id);

-- Episodes (only for series)
CREATE TABLE IF NOT EXISTS media_episode (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  season_id UUID NOT NULL REFERENCES media_season(id) ON DELETE CASCADE,
  episode_number INT NOT NULL,
  title TEXT,
  air_date DATE,
  tmdb_id BIGINT,
  tvdb_id BIGINT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (season_id, episode_number)
);

CREATE INDEX IF NOT EXISTS idx_media_episode_season ON media_episode (season_id);
CREATE INDEX IF NOT EXISTS idx_media_episode_tmdb ON media_episode (tmdb_id);
CREATE INDEX IF NOT EXISTS idx_media_episode_tvdb ON media_episode (tvdb_id);

-- Physical media files (movie or per-episode)
-- Removed season_id (derive from episode_id) and status (moved to media_file_state)
CREATE TABLE IF NOT EXISTS media_file (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  library_id UUID NOT NULL REFERENCES library(id) ON DELETE CASCADE,
  media_item_id UUID NOT NULL REFERENCES media_item(id) ON DELETE CASCADE,
  episode_id UUID REFERENCES media_episode(id) ON DELETE SET NULL,
  path TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CHECK (path !~ '^/'),
  UNIQUE (library_id, path)
);

CREATE INDEX IF NOT EXISTS idx_media_file_library ON media_file (library_id);
CREATE INDEX IF NOT EXISTS idx_media_file_media ON media_file (media_item_id);
CREATE INDEX IF NOT EXISTS idx_media_file_episode ON media_file (episode_id);

-- File state tracking (existence, size, verification)
CREATE TABLE IF NOT EXISTS media_file_state (
  media_file_id UUID PRIMARY KEY REFERENCES media_file(id) ON DELETE CASCADE,
  file_exists BOOLEAN NOT NULL DEFAULT true,
  file_size BIGINT,
  last_verified_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_media_file_state_exists ON media_file_state (file_exists);
CREATE INDEX IF NOT EXISTS idx_media_file_state_verified ON media_file_state (last_verified_at);

-- Import history (successes and failures)
CREATE TABLE IF NOT EXISTS media_file_import (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  media_file_id UUID REFERENCES media_file(id) ON DELETE SET NULL,
  download_job_id UUID,  -- FK added after download_job table is created
  method TEXT NOT NULL CHECK (method IN ('hardlink','copy','scan','manual_match')),
  source_path TEXT,
  dest_path TEXT NOT NULL,
  attempted_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  success BOOLEAN NOT NULL,
  error_message TEXT
);

CREATE INDEX IF NOT EXISTS idx_media_file_import_media_file ON media_file_import (media_file_id);
CREATE INDEX IF NOT EXISTS idx_media_file_import_download_job ON media_file_import (download_job_id);
CREATE INDEX IF NOT EXISTS idx_media_file_import_attempted ON media_file_import (attempted_at DESC);
CREATE INDEX IF NOT EXISTS idx_media_file_import_success ON media_file_import (success);

-- Unmatched files queue for manual matching
CREATE TABLE IF NOT EXISTS unmatched_file (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  library_id UUID NOT NULL REFERENCES library(id) ON DELETE CASCADE,
  path TEXT NOT NULL,
  file_size BIGINT,
  discovered_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  suggested_matches JSONB,
  resolved_at TIMESTAMPTZ,
  resolved_media_file_id UUID REFERENCES media_file(id) ON DELETE SET NULL,
  UNIQUE (library_id, path)
);

CREATE INDEX IF NOT EXISTS idx_unmatched_file_library ON unmatched_file (library_id);
CREATE INDEX IF NOT EXISTS idx_unmatched_file_resolved ON unmatched_file (resolved_at);
CREATE INDEX IF NOT EXISTS idx_unmatched_file_unresolved ON unmatched_file (library_id) WHERE resolved_at IS NULL;
