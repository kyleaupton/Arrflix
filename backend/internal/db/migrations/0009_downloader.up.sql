-- Downloaders store configuration for torrent/usenet clients
create table if not exists downloader (
  id uuid primary key default uuid_generate_v4(),
  name text not null,
  type text not null check (type in ('qbittorrent')),
  protocol text not null check (protocol in ('torrent', 'usenet')),
  url text not null,
  username text,
  password text,
  config_json jsonb,
  enabled boolean not null default true,
  "default" boolean not null default false,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create unique index if not exists uq_downloader_name_ci on downloader (lower(name));
create index if not exists idx_downloader_type on downloader (type);
create index if not exists idx_downloader_protocol on downloader (protocol);
create index if not exists idx_downloader_enabled on downloader (enabled);
create unique index if not exists uq_downloader_default_protocol on downloader (protocol) where "default" = true;

