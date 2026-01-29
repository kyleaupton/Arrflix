-- Download jobs track enqueue -> download -> import lifecycle (torrent/usenet)
-- Removed season_id for consistency (derive from episode_id)
CREATE TABLE IF NOT EXISTS download_job (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

  status TEXT NOT NULL CHECK (status IN (
    'created',
    'enqueued',
    'downloading',
    'completed',
    'importing',
    'imported',
    'failed',
    'cancelled'
  )),

  protocol TEXT NOT NULL CHECK (protocol IN ('torrent','usenet')),
  media_type TEXT NOT NULL CHECK (media_type IN ('movie','series')),

  -- Internal media graph target
  media_item_id UUID REFERENCES media_item(id) ON DELETE RESTRICT,
  episode_id UUID REFERENCES media_episode(id) ON DELETE RESTRICT,

  -- What the user picked (from indexer cache)
  indexer_id BIGINT NOT NULL,
  guid TEXT NOT NULL,
  candidate_title TEXT NOT NULL,
  candidate_link TEXT NOT NULL,

  -- Policy plan snapshot
  downloader_id UUID NOT NULL REFERENCES downloader(id) ON DELETE RESTRICT,
  library_id UUID NOT NULL REFERENCES library(id) ON DELETE RESTRICT,
  name_template_id UUID NOT NULL REFERENCES name_template(id) ON DELETE RESTRICT,

  -- Downloader linkage
  downloader_external_id TEXT,

  -- Best-effort download locations as seen by the downloader
  download_save_path TEXT,
  download_content_path TEXT,

  -- Import decisions/results
  import_source_path TEXT,
  import_dest_path TEXT,
  import_method TEXT CHECK (import_method IN ('hardlink','copy')),
  primary_media_file_id UUID REFERENCES media_file(id) ON DELETE SET NULL,

  -- Pre-calculated destination path with .{ext} placeholder
  predicted_dest_path TEXT,

  -- Progress snapshots (optional)
  downloader_status TEXT,
  progress DOUBLE PRECISION,

  -- Retry/backoff
  attempt_count INT NOT NULL DEFAULT 0,
  next_run_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  last_error TEXT,

  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_download_job_status ON download_job (status);
CREATE INDEX IF NOT EXISTS idx_download_job_next_run ON download_job (next_run_at);
CREATE INDEX IF NOT EXISTS idx_download_job_created_at ON download_job (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_download_job_media_item ON download_job (media_item_id);
CREATE INDEX IF NOT EXISTS idx_download_job_episode ON download_job (episode_id);
CREATE UNIQUE INDEX IF NOT EXISTS uq_download_job_candidate
  ON download_job (indexer_id, guid);

-- Link download jobs to media files (for multi-file downloads like season packs)
CREATE TABLE IF NOT EXISTS download_job_media_file (
  download_job_id UUID NOT NULL REFERENCES download_job(id) ON DELETE CASCADE,
  media_file_id UUID NOT NULL REFERENCES media_file(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (download_job_id, media_file_id)
);

CREATE INDEX IF NOT EXISTS idx_djmf_media_file ON download_job_media_file (media_file_id);

-- Add FK from media_file_import to download_job now that download_job exists
ALTER TABLE media_file_import
  ADD CONSTRAINT fk_media_file_import_download_job
  FOREIGN KEY (download_job_id)
  REFERENCES download_job(id)
  ON DELETE SET NULL;
