SELECT
  id,
  name,
  public_id,
  description,
  is_default,
  created_at,
  updated_at
FROM
  views
WHERE
  deleted_at IS NULL
  AND (
    @search::TEXT = ''
    OR name ILIKE '%' || @search || '%'
  )
ORDER BY
  is_default DESC, name ASC
LIMIT @limit;
