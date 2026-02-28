INSERT INTO incident_updates (incident_id, description, provider_id, provider_status, status, status_time)
VALUES (@incident_id, @description, @provider_id, @provider_status, @status, @status_time)
ON CONFLICT (provider_id) DO UPDATE
SET
  description = EXCLUDED.description,
  provider_status = EXCLUDED.provider_status,
  status = EXCLUDED.status,
  status_time = EXCLUDED.status_time,
  updated_at = NOW()
RETURNING *
