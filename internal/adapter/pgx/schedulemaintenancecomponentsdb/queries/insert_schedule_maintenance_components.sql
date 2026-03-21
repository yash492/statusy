INSERT INTO schedule_maintenance_components (scheduled_maintenance_id, component_id)
VALUES (@scheduled_maintenance_id, @component_id)
ON CONFLICT (scheduled_maintenance_id, component_id) DO UPDATE
SET updated_at = NOW()
RETURNING *
