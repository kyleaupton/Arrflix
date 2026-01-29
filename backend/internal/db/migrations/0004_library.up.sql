-- Libraries store media roots by type
CREATE TABLE IF NOT EXISTS library (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name TEXT NOT NULL,
  type TEXT NOT NULL CHECK (type IN ('movie','series')),
  root_path TEXT NOT NULL,
  enabled BOOLEAN NOT NULL DEFAULT true,
  "default" BOOLEAN NOT NULL DEFAULT false,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_library_name_ci ON library (lower(name));
CREATE INDEX IF NOT EXISTS idx_library_enabled ON library (enabled);
CREATE INDEX IF NOT EXISTS idx_library_type ON library (type);
CREATE UNIQUE INDEX IF NOT EXISTS uq_library_default_type ON library (type) WHERE "default" = true;
