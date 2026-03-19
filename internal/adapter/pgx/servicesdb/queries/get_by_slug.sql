SELECT id, title, slug
FROM services
WHERE slug = @slug AND deleted_at IS NULL