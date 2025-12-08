-- Name templates define file naming patterns for movies and series
create table if not exists name_template (
  id uuid primary key default uuid_generate_v4(),
  name text not null,
  type text not null check (type in ('movie','series')),
  template text not null,
  "default" boolean not null default false,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create unique index if not exists uq_name_template_name_ci on name_template (lower(name));
create index if not exists idx_name_template_type on name_template (type);
create unique index if not exists uq_name_template_default_type on name_template (type) where "default" = true;

