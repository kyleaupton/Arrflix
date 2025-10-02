-- 0002_seed_roles_permissions.sql
insert into role (name, description, built_in) values
  ('admin','Full control',true),
  ('manager','Manage libraries and requests',true),
  ('user','Standard user',true),
  ('guest','Read-only access',true)
on conflict (name) do nothing;

insert into permission (key, description) values
  -- library
  ('library.read','View libraries and items'),
  ('library.write','Create/update/delete libraries'),
  ('library.scan','Trigger scans'),
  -- media
  ('media.read','View media metadata and files'),
  ('media.write','Edit media metadata / re-link files'),
  -- requests
  ('requests.create','Create content requests'),
  ('requests.approve','Approve/deny requests'),
  -- jobs
  ('jobs.read','View queue and job states'),
  ('jobs.manage','Retry/ban/cancel jobs'),
  -- admin
  ('admin.settings.read','View app settings'),
  ('admin.settings.write','Change app settings'),
  ('admin.users.manage','Invite/disable users, manage roles')
on conflict (key) do nothing;

-- Grant sensible defaults to built-in roles (global scope)
-- Admin: everything
insert into permission_grant (subject_type, subject_id, permission_key, effect)
select 'role', r.id, p.key, 'allow'
from role r, permission p
where r.name='admin'
on conflict do nothing;

-- Manager: operational control, not user admin
insert into permission_grant (subject_type, subject_id, permission_key, effect)
select 'role', r.id, p.key, 'allow'
from role r
join permission p on p.key in (
  'library.read','library.write','library.scan',
  'media.read','media.write',
  'requests.create','requests.approve',
  'jobs.read','jobs.manage',
  'admin.settings.read'
)
where r.name='manager'
on conflict do nothing;

-- User: read + request create
insert into permission_grant (subject_type, subject_id, permission_key, effect)
select 'role', r.id, p.key, 'allow'
from role r
join permission p on p.key in (
  'library.read','media.read','requests.create','jobs.read'
)
where r.name='user'
on conflict do nothing;

-- Guest: read-only library + media
insert into permission_grant (subject_type, subject_id, permission_key, effect)
select 'role', r.id, p.key, 'allow'
from role r
join permission p on p.key in ('library.read','media.read')
where r.name='guest'
on conflict do nothing;
