INSERT INTO
  services (
    name,
    slug,
    incidents_url,
    schedule_maintenances_url,
    components_url,
    provider_type
  )
VALUES
  (
    @name,
    @slug,
    @incidents_url,
    @schedule_maintenances_url,
    @components_url,
    @provider_type
  )
ON CONFLICT (slug) DO UPDATE
SET
  name = EXCLUDED.name,
  incidents_url = EXCLUDED.incidents_url,
  schedule_maintenances_url = EXCLUDED.schedule_maintenances_url,
  components_url = EXCLUDED.components_url,
  provider_type = EXCLUDED.provider_type
RETURNING
  *;