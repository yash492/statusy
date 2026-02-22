INSERT INTO
  component_groups (
    name,
    provider_id,
    service_id
  )
VALUES
  (
    @name,
    @provider_id,
    @service_id
  )
ON CONFLICT (provider_id) DO UPDATE
SET
  name = EXCLUDED.name
RETURNING
  *;