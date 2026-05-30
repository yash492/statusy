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
  is_default = TRUE
  AND deleted_at IS NULL
ORDER BY
  id ASC
LIMIT
  1;
