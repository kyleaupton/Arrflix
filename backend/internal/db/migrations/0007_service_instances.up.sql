-- Service instances for dynamic container orchestration
create table service_instance (
  id uuid primary key default uuid_generate_v4(),
  name text not null unique,
  type text not null, -- 'qbittorrent', 'transmission', etc.
  enabled boolean not null default true,
  config jsonb not null default '{}'::jsonb,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create index idx_service_instance_type on service_instance(type);
create index idx_service_instance_enabled on service_instance(enabled);
