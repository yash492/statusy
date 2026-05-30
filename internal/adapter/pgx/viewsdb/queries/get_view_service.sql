SELECT
  id,
  view_id,
  service_id,
  include_all_components,
  created_at,
  updated_at
FROM
  view_services
WHERE
  view_id = @view_id
  AND service_id = @service_id
  AND deleted_at IS NULL
LIMIT 1;
