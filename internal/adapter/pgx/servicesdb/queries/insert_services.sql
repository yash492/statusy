INSERT INTO
  services (
    name,
    slug,
    url
  )
VALUES
  (
    @name,
    @slug,
    @url
  )
ON CONFLICT (slug) DO UPDATE
SET
  name = EXCLUDED.name,
  url = EXCLUDED.url
RETURNING
  id, name, slug, url;