-- Key/Value application settings with JSON values
CREATE TABLE IF NOT EXISTS app_setting (
  key TEXT PRIMARY KEY,
  type TEXT NOT NULL,
  value_json JSONB NOT NULL,
  version INT NOT NULL DEFAULT 1,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_app_setting_updated_at ON app_setting (updated_at DESC);

-- Track system initialization status
INSERT INTO app_setting (key, type, value_json)
VALUES ('system.initialized', 'bool', 'false'::jsonb)
ON CONFLICT (key) DO NOTHING;

CREATE INDEX IF NOT EXISTS idx_app_setting_system_initialized
  ON app_setting (key)
  WHERE key = 'system.initialized';
