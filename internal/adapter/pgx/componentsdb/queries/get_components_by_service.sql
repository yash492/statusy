SELECT
  id,
  name,
  provider_id,
  service_id,
  component_group_id,
  created_at,
  updated_at,
  deleted_at
FROM
  components
WHERE
  service_id = @service_id
  AND deleted_at IS NULL
ORDER BY
  name ASC;
