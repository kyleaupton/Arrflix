-- Track system initialization status using existing app_setting table
INSERT INTO app_setting (key, type, value_json)
VALUES ('system.initialized', 'bool', 'false'::jsonb)
ON CONFLICT (key) DO NOTHING;

-- Add index for fast initialization status checks
CREATE INDEX IF NOT EXISTS idx_app_setting_system_initialized
  ON app_setting (key)
  WHERE key = 'system.initialized';
