INSERT INTO
  components (
    name,
    provider_id,
    service_id,
    component_group_id
  )
VALUES
  (
    @name,
    @provider_id,
    @service_id,
    @component_group_id
  )
ON CONFLICT (provider_id) DO UPDATE
SET
  name = EXCLUDED.name
RETURNING
  *;