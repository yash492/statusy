UPDATE views
SET
  name = @name,
  slug = @slug,
  description = @description,
  is_default = @is_default,
  updated_at = now()
WHERE
  id = @id
  AND deleted_at IS NULL
RETURNING
  *;
