UPDATE view_services
SET
  deleted_at = now()
WHERE
  view_id = @view_id
  AND service_id = @service_id
  AND deleted_at IS NULL
RETURNING
  id;
