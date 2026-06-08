UPDATE view_notifications
SET name = $2, type = $3, config = $4, updated_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING id, view_id, name, type, config, created_at, updated_at;
