INSERT INTO
  views (
    name,
    public_id,
    description,
    is_default
  )
VALUES
  (
    @name,
    @public_id,
    @description,
    @is_default
  )
RETURNING
  *;