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
  public_id = @public_id
  AND deleted_at IS NULL
LIMIT 1;
