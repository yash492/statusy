INSERT INTO
  services (
    name,
    slug
  )
VALUES
  (
    @name,
    @slug
  )
ON CONFLICT (slug) DO UPDATE
SET
  name = EXCLUDED.name
RETURNING
  id, name, slug;