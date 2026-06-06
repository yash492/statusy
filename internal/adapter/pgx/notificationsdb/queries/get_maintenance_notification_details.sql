SELECT 
	sm.id AS maintenance_id,
	smu.id AS update_id,
	sm.title,
	smu.status,
	smu.description,
	sm.provider_id,
	s.name AS service_name,
	COALESCE(array_agg(c.name) FILTER (WHERE c.name IS NOT NULL), '{}') AS component_names,
	sm.start_time,
	sm.end_time,
	smu.updated_at
FROM scheduled_maintenance_updates smu
JOIN scheduled_maintenances sm ON sm.id = smu.scheduled_maintenance_id
JOIN services s ON s.id = sm.service_id
LEFT JOIN scheduled_maintenance_components smc ON smc.scheduled_maintenance_id = sm.id
LEFT JOIN components c ON c.id = smc.component_id
WHERE smu.id = $1
GROUP BY sm.id, smu.id, sm.title, smu.status, smu.description, sm.provider_id, s.name, sm.start_time, sm.end_time, smu.updated_at;
