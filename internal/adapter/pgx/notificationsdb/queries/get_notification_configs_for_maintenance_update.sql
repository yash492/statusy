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
JOIN scheduled_maintenance_updates smu ON smu.id = $1
JOIN scheduled_maintenances sm ON sm.id = smu.scheduled_maintenance_id AND sm.service_id = vs.service_id
WHERE vs.monitor_scheduled_maintenances = true
  AND vn.deleted_at IS NULL
  AND (
	vs.include_all_components = true
	OR NOT EXISTS (SELECT 1 FROM scheduled_maintenance_components smc WHERE smc.scheduled_maintenance_id = sm.id)
	OR EXISTS (
	  SELECT 1 FROM scheduled_maintenance_components smc
	  JOIN view_service_components vsc ON vsc.component_id = smc.component_id AND vsc.deleted_at IS NULL
	  WHERE smc.scheduled_maintenance_id = sm.id AND vsc.view_service_id = vs.id
	)
	OR EXISTS (
	  SELECT 1 FROM scheduled_maintenance_components smc
	  JOIN components c ON c.id = smc.component_id AND c.deleted_at IS NULL
	  JOIN view_service_component_groups vscg ON vscg.component_group_id = c.component_group_id AND vscg.deleted_at IS NULL
	  WHERE smc.scheduled_maintenance_id = sm.id AND vscg.view_service_id = vs.id
	)
  );
