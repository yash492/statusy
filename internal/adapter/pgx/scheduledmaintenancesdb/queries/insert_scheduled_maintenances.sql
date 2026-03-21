INSERT INTO scheduled_maintenances (
  title,
  link,
  provider_impact,
  impact,
  starts_at,
  ends_at,
  service_id,
  provider_id,
  provider_created_at  
) VALUES (
  @title,
  @link,
  @provider_impact,
  @impact,
  @starts_at,
  @ends_at,
  @service_id,
  @provider_id,
  @provider_created_at
)
ON CONFLICT (provider_id) DO UPDATE
SET
  title = EXCLUDED.title,
  link = EXCLUDED.link,
  provider_impact = EXCLUDED.provider_impact,
  impact = EXCLUDED.impact,
  updated_at = NOW()
RETURNING *;