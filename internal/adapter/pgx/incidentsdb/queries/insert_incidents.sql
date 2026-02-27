INSERT INTO incidents (
  name,
  link,
  provider_impact,
  impact,
  service_id,
  provider_id,
  provider_created_at,
  created_at,
  updated_at
) VALUES (
  @name,
  @link,
  @provider_impact,
  @impact,
  @service_id,
  @provider_id,
  @provider_created_at
)
ON CONFLICT DO NOTHING
RETURNING *;