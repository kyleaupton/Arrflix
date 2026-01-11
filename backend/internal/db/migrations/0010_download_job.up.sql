-- Download jobs track enqueue -> download -> import lifecycle (torrent/usenet)
create table if not exists download_job (
  id uuid primary key default uuid_generate_v4(),

  status text not null check (status in (
    'created',
    'enqueued',
    'downloading',
    'completed',
    'importing',
    'imported',
    'failed',
    'cancelled'
  )),

  protocol text not null check (protocol in ('torrent','usenet')),
  media_type text not null check (media_type in ('movie','series')),

  -- Internal media graph target
  media_item_id uuid references media_item(id) on delete restrict,
  season_id uuid references media_season(id) on delete restrict,
  episode_id uuid references media_episode(id) on delete restrict,

  -- What the user picked (from indexer cache)
  indexer_id bigint not null,
  guid text not null,
  candidate_title text not null,
  candidate_link text not null,

  -- Policy plan snapshot
  downloader_id uuid not null references downloader(id) on delete restrict,
  library_id uuid not null references library(id) on delete restrict,
  name_template_id uuid not null references name_template(id) on delete restrict,

  -- Downloader linkage
  downloader_external_id text,

  -- Best-effort download locations as seen by the downloader
  download_save_path text,
  download_content_path text,

  -- Import decisions/results
  import_source_path text,
  import_dest_path text,
  import_method text check (import_method in ('hardlink','copy')),
  primary_media_file_id uuid references media_file(id) on delete set null,

  -- Progress snapshots (optional)
  downloader_status text,
  progress double precision,

  -- Retry/backoff
  attempt_count int not null default 0,
  next_run_at timestamptz not null default now(),
  last_error text,

  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create index if not exists idx_download_job_status on download_job (status);
create index if not exists idx_download_job_next_run on download_job (next_run_at);
create index if not exists idx_download_job_created_at on download_job (created_at desc);
create index if not exists idx_download_job_media_item on download_job (media_item_id);
create index if not exists idx_download_job_episode on download_job (episode_id);
create unique index if not exists uq_download_job_candidate
  on download_job (indexer_id, guid);

create table if not exists download_job_media_file (
  download_job_id uuid not null references download_job(id) on delete cascade,
  media_file_id uuid not null references media_file(id) on delete cascade,
  created_at timestamptz not null default now(),
  unique (download_job_id, media_file_id)
);

create index if not exists idx_djmf_media_file on download_job_media_file (media_file_id);


