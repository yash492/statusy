UPDATE view_services
SET
  include_all_components = @include_all_components,
  updated_at = now()
WHERE
  id = @id
  AND deleted_at IS NULL
RETURNING
  id,
  view_id,
  service_id,
  include_all_components,
  created_at,
  updated_at;
