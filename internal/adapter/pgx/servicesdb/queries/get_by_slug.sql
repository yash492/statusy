SELECT id, name, slug
FROM services
WHERE slug = @slug AND deleted_at IS NULL