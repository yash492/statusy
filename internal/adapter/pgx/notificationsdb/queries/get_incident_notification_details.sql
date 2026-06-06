SELECT 
	i.id AS incident_id,
	iu.id AS update_id,
	i.title,
	iu.status,
	iu.description,
	i.provider_id,
	s.name AS service_name,
	COALESCE(array_agg(c.name) FILTER (WHERE c.name IS NOT NULL), '{}') AS component_names,
	iu.updated_at
FROM incident_updates iu
JOIN incidents i ON i.id = iu.incident_id
JOIN services s ON s.id = i.service_id
LEFT JOIN incident_components ic ON ic.incident_id = i.id
LEFT JOIN components c ON c.id = ic.component_id
WHERE iu.id = $1
GROUP BY i.id, iu.id, i.title, iu.status, iu.description, i.provider_id, s.name, iu.updated_at;
