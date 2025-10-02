create extension if not exists "uuid-ossp";

-- Users
create table if not exists app_user (
  id uuid primary key default uuid_generate_v4(),
  email text unique,                    -- nullable for pure Plex users
  display_name text,
  avatar_url text,
  password_hash text,                   -- null when SSO-only
  is_active boolean not null default true,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create index if not exists idx_app_user_email_ci
  on app_user (lower(email));

-- External identities (Plex, etc.)
-- A user may have multiple identities; identities are unique across provider+subject.
create type auth_provider as enum ('local','plex');

create table if not exists user_identity (
  id uuid primary key default uuid_generate_v4(),
  user_id uuid not null references app_user(id) on delete cascade,
  provider auth_provider not null,
  subject text not null,                -- external user id / subject
  username text,                        -- external username/handle
  access_token text,                    -- if you store Plex token (encrypt/secret-manage in app)
  refresh_token text,
  token_expires_at timestamptz,
  raw jsonb not null default '{}'::jsonb,  -- raw profile snapshot
  created_at timestamptz not null default now(),
  unique (provider, subject)
);

create index if not exists idx_identity_user ON user_identity (user_id);

-- Roles (RBAC)
create table if not exists role (
  id uuid primary key default uuid_generate_v4(),
  name text not null unique,            -- e.g., 'admin','manager','user','guest'
  description text,
  built_in boolean not null default false,
  created_at timestamptz not null default now()
);

-- User <-> Role (membership)
create table if not exists user_role (
  user_id uuid not null references app_user(id) on delete cascade,
  role_id uuid not null references role(id) on delete cascade,
  granted_at timestamptz not null default now(),
  primary key (user_id, role_id)
);

-- Permissions catalog (flat keys; hierarchical by naming convention)
-- Examples: 'library.read', 'library.write', 'requests.create', 'requests.approve', 'admin.*'
create table if not exists permission (
  id uuid primary key default uuid_generate_v4(),
  key text not null unique,             -- unique machine key
  description text
);

-- Unified grants (fine-grained + role grants + user overrides)
-- A grant applies to either a role OR a user (subject_type + subject_id).
-- It can be global (no resource) or scoped to a resource (type+id).
-- 'effect' supports explicit deny to override an allow.
create type grant_subject as enum ('role','user');
create type grant_effect as enum ('allow','deny');

-- NOTE: use table name permission_grant to avoid keyword conflicts with GRANT
create table if not exists permission_grant (
  id uuid primary key default uuid_generate_v4(),
  subject_type grant_subject not null,
  subject_id uuid not null,             -- role.id or app_user.id depending on subject_type
  permission_key text not null references permission(key) on delete cascade,
  resource_type text,                   -- e.g., 'library'
  resource_id uuid,                     -- id of a specific resource (nullable for global)
  effect grant_effect not null default 'allow',
  created_at timestamptz not null default now()
  -- Prevent duplicates; allow different scope or different effects
);

create index if not exists idx_grant_subject ON permission_grant (subject_type, subject_id);
create index if not exists idx_grant_perm ON permission_grant (permission_key);
create index if not exists idx_grant_scope ON permission_grant (resource_type, resource_id);

-- Enforce uniqueness across nullable scope using an expression unique index
create unique index if not exists uq_permission_grant_natural
  on permission_grant (
    subject_type,
    subject_id,
    permission_key,
    coalesce(resource_type,'*'),
    coalesce(resource_id::text,'00000000-0000-0000-0000-000000000000')
  );

-- (Optional) basic audit trail (you can expand later)
create table if not exists auth_audit (
  id uuid primary key default uuid_generate_v4(),
  user_id uuid references app_user(id) on delete set null,
  event text not null,                  -- 'login.success','login.failure','token.refresh','role.assign', ...
  detail jsonb not null default '{}'::jsonb,
  created_at timestamptz not null default now()
);
