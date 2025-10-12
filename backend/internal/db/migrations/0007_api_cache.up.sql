create table if not exists api_cache (
	id           uuid primary key default gen_random_uuid(),
	key          text not null unique,
	category     text,
	response     jsonb not null,
	status       integer not null default 200,
	content_type text,
	headers      jsonb,
	stored_at    timestamptz not null default now(),
	expires_at   timestamptz not null
);

create index if not exists idx_api_cache_expires_at on api_cache (expires_at);
create index if not exists idx_api_cache_category on api_cache (category);

-- You will typically query by key and check expires_at > now(), so a btree on key is fine.
create index if not exists idx_api_cache_key on api_cache (key);