-- Policies define rules and actions for torrent processing
create table if not exists policy (
  id uuid primary key default uuid_generate_v4(),
  name text not null,
  description text,
  enabled boolean not null default true,
  priority integer not null,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create index if not exists idx_policy_priority on policy (priority desc);
create index if not exists idx_policy_enabled on policy (enabled);

-- Rules define conditions for policies
create table if not exists rule (
  id uuid primary key default uuid_generate_v4(),
  policy_id uuid not null references policy(id) on delete cascade,
  left_operand text not null,
  operator text not null check (operator in ('==', '!=', '>', '>=', '<', '<=', 'contains', 'in', 'not in', 'and', 'or', 'not')),
  right_operand text not null,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  unique (policy_id)
);

create index if not exists idx_rule_policy on rule (policy_id);

-- Actions define mutations to apply when policy matches
create table if not exists action (
  id uuid primary key default uuid_generate_v4(),
  policy_id uuid not null references policy(id) on delete cascade,
  type text not null check (type in ('set_downloader', 'set_library', 'set_name_template', 'stop_processing')),
  value text not null,
  "order" integer not null,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create index if not exists idx_action_policy_order on action (policy_id, "order");

