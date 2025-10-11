-- Libraries store media roots by type
create table if not exists library (
  id uuid primary key default uuid_generate_v4(),
  name text not null,
  type text not null check (type in ('movie','series')),
  root_path text not null,
  enabled boolean not null default true,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create unique index if not exists uq_library_name_ci on library (lower(name));
create index if not exists idx_library_enabled on library (enabled);
create index if not exists idx_library_type on library (type);


