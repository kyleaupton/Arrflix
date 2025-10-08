-- Media core schema: media_item, media_season, media_episode, media_file

-- Media items (movies or series)
create table if not exists media_item (
  id uuid primary key default uuid_generate_v4(),
  library_id uuid not null references library(id) on delete cascade,
  type text not null check (type in ('movie','series')),
  title text not null,
  year int,
  tmdb_id int,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create index if not exists idx_media_item_library on media_item (library_id);
create index if not exists idx_media_item_type on media_item (type);
create index if not exists idx_media_item_tmdb on media_item (tmdb_id);

-- Seasons (only for series)
create table if not exists media_season (
  id uuid primary key default uuid_generate_v4(),
  media_id uuid not null references media_item(id) on delete cascade,
  season_number int not null,
  title text,
  air_date date,
  created_at timestamptz not null default now(),
  unique (media_id, season_number)
);

create index if not exists idx_media_season_media on media_season (media_id);

-- Episodes (only for series)
create table if not exists media_episode (
  id uuid primary key default uuid_generate_v4(),
  season_id uuid not null references media_season(id) on delete cascade,
  episode_number int not null,
  title text,
  air_date date,
  tmdb_id bigint,
  tvdb_id bigint,
  created_at timestamptz not null default now(),
  unique (season_id, episode_number)
);

create index if not exists idx_media_episode_season on media_episode (season_id);
create index if not exists idx_media_episode_tmdb on media_episode (tmdb_id);
create index if not exists idx_media_episode_tvdb on media_episode (tvdb_id);

-- Physical media files (movie or per-episode)
create table if not exists media_file (
  id uuid primary key default uuid_generate_v4(),
  media_id uuid not null references media_item(id) on delete cascade,
  season_id uuid references media_season(id) on delete set null,
  episode_id uuid references media_episode(id) on delete set null,
  path text not null unique,
  size_bytes bigint,
  resolution text,
  added_at timestamptz not null default now()
);

create index if not exists idx_media_file_media on media_file (media_id);
create index if not exists idx_media_file_episode on media_file (episode_id);




