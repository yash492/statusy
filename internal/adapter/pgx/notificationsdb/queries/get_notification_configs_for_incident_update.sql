SELECT DISTINCT
	vn.id,
	vn.view_id,
	vn.type,
	vn.config,
	vn.created_at,
	vn.updated_at
FROM view_notifications vn
JOIN views v ON v.id = vn.view_id AND v.deleted_at IS NULL
JOIN view_services vs ON vs.view_id = v.id AND vs.deleted_at IS NULL
JOIN incident_updates iu ON iu.id = $1
JOIN incidents i ON i.id = iu.incident_id AND i.service_id = vs.service_id
WHERE vs.monitor_incidents = true
  AND vn.deleted_at IS NULL
  AND (
	vs.include_all_components = true
	OR NOT EXISTS (SELECT 1 FROM incident_components ic WHERE ic.incident_id = i.id)
	OR EXISTS (
	  SELECT 1 FROM incident_components ic
	  JOIN view_service_components vsc ON vsc.component_id = ic.component_id AND vsc.deleted_at IS NULL
	  WHERE ic.incident_id = i.id AND vsc.view_service_id = vs.id
	)
	OR EXISTS (
	  SELECT 1 FROM incident_components ic
	  JOIN components c ON c.id = ic.component_id AND c.deleted_at IS NULL
	  JOIN view_service_component_groups vscg ON vscg.component_group_id = c.component_group_id AND vscg.deleted_at IS NULL
	  WHERE ic.incident_id = i.id AND vscg.view_service_id = vs.id
	)
  );
