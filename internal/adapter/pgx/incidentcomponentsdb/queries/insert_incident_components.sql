INSERT INTO incident_components (incident_id, component_id)
VALUES (@incident_id, @component_id)
ON CONFLICT (incident_id, component_id) DO UPDATE
SET updated_at = NOW()
RETURNING *
