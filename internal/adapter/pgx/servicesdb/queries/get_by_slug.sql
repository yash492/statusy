SELECT id, name, slug, url
FROM services
WHERE slug = @slug AND deleted_at IS NULL