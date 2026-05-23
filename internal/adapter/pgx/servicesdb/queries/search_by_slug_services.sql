SELECT id, name, slug
FROM services
WHERE name ILIKE @name
	AND deleted_at IS NULL;