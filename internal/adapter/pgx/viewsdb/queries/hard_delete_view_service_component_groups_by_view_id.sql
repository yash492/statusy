DELETE FROM view_service_component_groups
WHERE view_service_id IN (
  SELECT id FROM view_services WHERE view_id = @view_id AND deleted_at IS NULL
) AND deleted_at IS NULL;
