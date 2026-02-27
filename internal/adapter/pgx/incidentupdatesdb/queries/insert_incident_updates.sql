INSERT INTO incident_updates (incident_id, description, provider_id, provider_status, status, status_time)
VALUES (@incident_id, @description, @provider_id, @provider_status, @status, @status_time)
ON CONFLICT DO NOTHING
RETURNING *
