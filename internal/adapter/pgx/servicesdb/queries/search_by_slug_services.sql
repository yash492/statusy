SELECT id, name, slug, url
FROM services
WHERE name ILIKE @name
	AND deleted_at IS NULL;