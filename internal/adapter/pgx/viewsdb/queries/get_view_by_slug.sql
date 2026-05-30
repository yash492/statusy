SELECT
  id,
  name,
  slug,
  description,
  is_default,
  created_at,
  updated_at
FROM
  views
WHERE
  slug = @slug
  AND deleted_at IS NULL
LIMIT 1;
