-- Seed roles and permissions
INSERT INTO role (name, description, built_in) VALUES
  ('admin','Full control',true),
  ('manager','Manage libraries and requests',true),
  ('user','Standard user',true),
  ('guest','Read-only access',true)
ON CONFLICT (name) DO NOTHING;

INSERT INTO permission (key, description) VALUES
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
ON CONFLICT (key) DO NOTHING;

-- Grant sensible defaults to built-in roles (global scope)
-- Admin: everything
INSERT INTO permission_grant (subject_type, subject_id, permission_key, effect)
SELECT 'role', r.id, p.key, 'allow'
FROM role r, permission p
WHERE r.name='admin'
ON CONFLICT DO NOTHING;

-- Manager: operational control, not user admin
INSERT INTO permission_grant (subject_type, subject_id, permission_key, effect)
SELECT 'role', r.id, p.key, 'allow'
FROM role r
JOIN permission p ON p.key IN (
  'library.read','library.write','library.scan',
  'media.read','media.write',
  'requests.create','requests.approve',
  'jobs.read','jobs.manage',
  'admin.settings.read'
)
WHERE r.name='manager'
ON CONFLICT DO NOTHING;

-- User: read + request create
INSERT INTO permission_grant (subject_type, subject_id, permission_key, effect)
SELECT 'role', r.id, p.key, 'allow'
FROM role r
JOIN permission p ON p.key IN (
  'library.read','media.read','requests.create','jobs.read'
)
WHERE r.name='user'
ON CONFLICT DO NOTHING;

-- Guest: read-only library + media
INSERT INTO permission_grant (subject_type, subject_id, permission_key, effect)
SELECT 'role', r.id, p.key, 'allow'
FROM role r
JOIN permission p ON p.key IN ('library.read','media.read')
WHERE r.name='guest'
ON CONFLICT DO NOTHING;
