INSERT INTO incidents (
  name,
  link,
  provider_impact,
  impact,
  service_id,
  provider_id,
  provider_created_at  
) VALUES (
  @name,
  @link,
  @provider_impact,
  @impact,
  @service_id,
  @provider_id,
  @provider_created_at
)
ON CONFLICT (provider_id) DO UPDATE
SET
  name = EXCLUDED.name,
  link = EXCLUDED.link,
  provider_impact = EXCLUDED.provider_impact,
  impact = EXCLUDED.impact,
  updated_at = NOW()
RETURNING *;