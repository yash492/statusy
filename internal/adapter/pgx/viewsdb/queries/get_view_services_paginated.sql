SELECT
  s.id AS id,
  s.name AS name,
  s.slug AS slug,
  vs.include_all_components AS include_all_components,
  vs.monitor_incidents AS monitor_incidents,
  vs.monitor_scheduled_maintenances AS monitor_scheduled_maintenances,
  CASE 
    WHEN vs.monitor_incidents = true AND ai.id IS NOT NULL THEN 'down'
    ELSE 'up'
  END AS status,
  COALESCE(ai.title, '') AS last_incident,
  COALESCE(ai.link, '') AS last_incident_link,
  COALESCE(um.title, '') AS upcoming_maintenance,
  COALESCE(um.link, '') AS upcoming_maintenance_link,
  COALESCE(vsc.component_ids, '{}') AS component_ids,
  COALESCE(vscg.component_group_ids, '{}') AS component_group_ids
FROM view_services vs
JOIN views v ON v.id = vs.view_id
JOIN services s ON s.id = vs.service_id
LEFT JOIN LATERAL (
  SELECT COALESCE(array_agg(component_id) FILTER (WHERE component_id IS NOT NULL), '{}') AS component_ids
  FROM view_service_components
  WHERE view_service_id = vs.id AND deleted_at IS NULL
) vsc ON true
LEFT JOIN LATERAL (
  SELECT COALESCE(array_agg(component_group_id) FILTER (WHERE component_group_id IS NOT NULL), '{}') AS component_group_ids
  FROM view_service_component_groups
  WHERE view_service_id = vs.id AND deleted_at IS NULL
) vscg ON true
LEFT JOIN LATERAL (
  SELECT i.id, i.title, i.link
  FROM incidents i
  WHERE i.service_id = s.id
    AND i.deleted_at IS NULL
    AND NOT i.is_resolved
    AND (
      vs.include_all_components = true
      OR NOT EXISTS (SELECT 1 FROM incident_components ic WHERE ic.incident_id = i.id)
      OR EXISTS (
        SELECT 1 FROM incident_components ic
        WHERE ic.incident_id = i.id
          AND (
            ic.component_id IN (SELECT component_id FROM view_service_components WHERE view_service_id = vs.id AND deleted_at IS NULL)
            OR ic.component_id IN (
              SELECT c.id FROM components c
              JOIN view_service_component_groups vscg ON c.component_group_id = vscg.component_group_id
              WHERE vscg.view_service_id = vs.id AND vscg.deleted_at IS NULL AND c.deleted_at IS NULL
            )
          )
      )
    )
  ORDER BY i.provider_created_at DESC, i.id DESC
  LIMIT 1
) ai ON true
LEFT JOIN LATERAL (
  SELECT sm.id, sm.title, sm.link
  FROM scheduled_maintenances sm
  WHERE sm.service_id = s.id
    AND sm.deleted_at IS NULL
    AND sm.ends_at > NOW()
    AND (
      vs.include_all_components = true
      OR NOT EXISTS (SELECT 1 FROM scheduled_maintenance_components smc WHERE smc.scheduled_maintenance_id = sm.id)
      OR EXISTS (
        SELECT 1 FROM scheduled_maintenance_components smc
        WHERE smc.scheduled_maintenance_id = sm.id
          AND (
            smc.component_id IN (SELECT component_id FROM view_service_components WHERE view_service_id = vs.id AND deleted_at IS NULL)
            OR smc.component_id IN (
              SELECT c.id FROM components c
              JOIN view_service_component_groups vscg ON c.component_group_id = vscg.component_group_id
              WHERE vscg.view_service_id = vs.id AND vscg.deleted_at IS NULL AND c.deleted_at IS NULL
            )
          )
      )
    )
  ORDER BY sm.starts_at ASC, sm.id DESC
  LIMIT 1
) um ON true
WHERE vs.view_id = @view_id
  AND vs.deleted_at IS NULL
  AND v.deleted_at IS NULL
  AND s.deleted_at IS NULL
  AND (
    @search::text = ''
    OR s.name ILIKE '%' || @search || '%'
    OR s.slug ILIKE '%' || @search || '%'
  )
ORDER BY 
  CASE 
    WHEN vs.monitor_incidents = true AND ai.id IS NOT NULL THEN 0 
    ELSE 1 
  END ASC,
  vs.id ASC
LIMIT @limit
OFFSET @offset;
