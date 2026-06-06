SELECT id, view_id, type, config, created_at, updated_at
FROM view_notifications
WHERE view_id = $1 AND deleted_at IS NULL;
