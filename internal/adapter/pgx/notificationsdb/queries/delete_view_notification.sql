UPDATE view_notifications
SET deleted_at = now(), updated_at = now()
WHERE id = $1;
