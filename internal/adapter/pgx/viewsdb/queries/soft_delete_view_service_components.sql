UPDATE view_service_components
SET
  deleted_at = now()
WHERE
  view_service_id = @view_service_id
  AND deleted_at IS NULL;
