INSERT INTO scheduled_maintenance_updates (scheduled_maintenance_id, description, provider_id, provider_status, status, status_time)
VALUES (@scheduled_maintenance_id, @description, @provider_id, @provider_status, @status, @status_time)
ON CONFLICT (provider_id) DO UPDATE
SET
  description = EXCLUDED.description,
  provider_status = EXCLUDED.provider_status,
  status = EXCLUDED.status,
  status_time = EXCLUDED.status_time,
  updated_at = NOW()
RETURNING *
