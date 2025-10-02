-- Key/Value application settings with JSON values
create table if not exists app_setting (
  key text primary key,
  type text not null,        -- text|bool|int|json
  value_json jsonb not null,
  version int not null default 1,
  updated_at timestamptz not null default now()
);

create index if not exists idx_app_setting_updated_at on app_setting (updated_at desc);


