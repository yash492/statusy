UPDATE view_service_components
SET deleted_at = now(), updated_at = now()
WHERE view_service_id IN (
  SELECT id FROM view_services WHERE view_id = @view_id AND deleted_at IS NULL
) AND deleted_at IS NULL;
