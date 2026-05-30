UPDATE views
SET deleted_at = now(), updated_at = now()
WHERE id = @view_id AND deleted_at IS NULL;
