-- Extensions and Authentication tables
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users
CREATE TABLE IF NOT EXISTS app_user (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  email TEXT UNIQUE,
  display_name TEXT,
  avatar_url TEXT,
  password_hash TEXT,
  is_active BOOLEAN NOT NULL DEFAULT true,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_app_user_email_ci
  ON app_user (lower(email));

-- External identities (Plex, etc.)
CREATE TYPE auth_provider AS ENUM ('local','plex');

CREATE TABLE IF NOT EXISTS user_identity (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL REFERENCES app_user(id) ON DELETE CASCADE,
  provider auth_provider NOT NULL,
  subject TEXT NOT NULL,
  username TEXT,
  access_token TEXT,
  refresh_token TEXT,
  token_expires_at TIMESTAMPTZ,
  raw JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (provider, subject)
);

CREATE INDEX IF NOT EXISTS idx_identity_user ON user_identity (user_id);

-- Roles (RBAC)
CREATE TABLE IF NOT EXISTS role (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name TEXT NOT NULL UNIQUE,
  description TEXT,
  built_in BOOLEAN NOT NULL DEFAULT false,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- User <-> Role (membership)
CREATE TABLE IF NOT EXISTS user_role (
  user_id UUID NOT NULL REFERENCES app_user(id) ON DELETE CASCADE,
  role_id UUID NOT NULL REFERENCES role(id) ON DELETE CASCADE,
  granted_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (user_id, role_id)
);

-- Permissions catalog
CREATE TABLE IF NOT EXISTS permission (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  key TEXT NOT NULL UNIQUE,
  description TEXT
);

-- Unified grants
CREATE TYPE grant_subject AS ENUM ('role','user');
CREATE TYPE grant_effect AS ENUM ('allow','deny');

CREATE TABLE IF NOT EXISTS permission_grant (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  subject_type grant_subject NOT NULL,
  subject_id UUID NOT NULL,
  permission_key TEXT NOT NULL REFERENCES permission(key) ON DELETE CASCADE,
  resource_type TEXT,
  resource_id UUID,
  effect grant_effect NOT NULL DEFAULT 'allow',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_grant_subject ON permission_grant (subject_type, subject_id);
CREATE INDEX IF NOT EXISTS idx_grant_perm ON permission_grant (permission_key);
CREATE INDEX IF NOT EXISTS idx_grant_scope ON permission_grant (resource_type, resource_id);

CREATE UNIQUE INDEX IF NOT EXISTS uq_permission_grant_natural
  ON permission_grant (
    subject_type,
    subject_id,
    permission_key,
    coalesce(resource_type,'*'),
    coalesce(resource_id::text,'00000000-0000-0000-0000-000000000000')
  );

-- Audit trail
CREATE TABLE IF NOT EXISTS auth_audit (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID REFERENCES app_user(id) ON DELETE SET NULL,
  event TEXT NOT NULL,
  detail JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
