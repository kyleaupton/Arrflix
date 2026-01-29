-- API response cache
CREATE TABLE IF NOT EXISTS api_cache (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  key TEXT NOT NULL UNIQUE,
  category TEXT,
  response JSONB NOT NULL,
  status INTEGER NOT NULL DEFAULT 200,
  content_type TEXT,
  headers JSONB,
  stored_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  expires_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_api_cache_expires_at ON api_cache (expires_at);
CREATE INDEX IF NOT EXISTS idx_api_cache_category ON api_cache (category);
CREATE INDEX IF NOT EXISTS idx_api_cache_key ON api_cache (key);
