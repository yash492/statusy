INSERT INTO
  services (
    title,
    slug
  )
VALUES
  (
    @title,
    @slug
  )
ON CONFLICT (slug) DO UPDATE
SET
  title = EXCLUDED.title
RETURNING
  id, title, slug;