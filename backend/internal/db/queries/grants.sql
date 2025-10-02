-- Returns grants that apply to the user directly or via any of the user's roles
SELECT g.*
FROM permission_grant g
WHERE
  (
    (g.subject_type = 'user' AND g.subject_id = $1)
    OR
    (g.subject_type = 'role' AND g.subject_id = ANY($2::uuid[]))
  )
  AND (g.permission_key = ANY($3::text[]))
  AND (
    (g.resource_type is null and g.resource_id is null) OR
    (g.resource_type = $4 and g.resource_id = $5)
  );