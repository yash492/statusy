INSERT INTO
  views (
    name,
    slug,
    description,
    is_default
  )
VALUES
  (
    @name,
    @slug,
    @description,
    @is_default
  )
RETURNING
  *;