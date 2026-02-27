INSERT INTO incident_components (incident_id, component_id)
VALUES (@incident_id, @component_id)
ON CONFLICT DO NOTHING
RETURNING *
