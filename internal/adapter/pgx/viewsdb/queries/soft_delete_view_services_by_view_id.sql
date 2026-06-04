UPDATE view_services
SET deleted_at = now(), updated_at = now()
WHERE view_id = @view_id AND deleted_at IS NULL;
