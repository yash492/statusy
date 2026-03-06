SELECT *
FROM services
WHERE slug ILIKE @slug
	AND deleted_at IS NULL;