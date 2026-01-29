-- Policies define rules and actions for torrent processing
CREATE TABLE IF NOT EXISTS policy (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name TEXT NOT NULL,
  description TEXT,
  enabled BOOLEAN NOT NULL DEFAULT true,
  priority INTEGER NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_policy_priority ON policy (priority DESC);
CREATE INDEX IF NOT EXISTS idx_policy_enabled ON policy (enabled);

-- Rules define conditions for policies
CREATE TABLE IF NOT EXISTS rule (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  policy_id UUID NOT NULL REFERENCES policy(id) ON DELETE CASCADE,
  left_operand TEXT NOT NULL,
  operator TEXT NOT NULL CHECK (operator IN ('==', '!=', '>', '>=', '<', '<=', 'contains', 'in', 'not in', 'and', 'or', 'not')),
  right_operand TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (policy_id)
);

CREATE INDEX IF NOT EXISTS idx_rule_policy ON rule (policy_id);

-- Actions define mutations to apply when policy matches
CREATE TABLE IF NOT EXISTS action (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  policy_id UUID NOT NULL REFERENCES policy(id) ON DELETE CASCADE,
  type TEXT NOT NULL CHECK (type IN ('set_downloader', 'set_library', 'set_name_template', 'stop_processing')),
  value TEXT NOT NULL,
  "order" INTEGER NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_action_policy_order ON action (policy_id, "order");

-- Name templates define file naming patterns for movies and series
CREATE TABLE IF NOT EXISTS name_template (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name TEXT NOT NULL,
  type TEXT NOT NULL CHECK (type IN ('movie','series')),
  template TEXT NOT NULL,
  movie_dir_template TEXT,
  series_show_template TEXT,
  series_season_template TEXT,
  "default" BOOLEAN NOT NULL DEFAULT false,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_name_template_name_ci ON name_template (lower(name));
CREATE INDEX IF NOT EXISTS idx_name_template_type ON name_template (type);
CREATE UNIQUE INDEX IF NOT EXISTS uq_name_template_default_type ON name_template (type) WHERE "default" = true;
