SELECT
  vs.id,
  vs.view_id,
  vs.service_id,
  vs.include_all_components,
  vs.created_at,
  vs.updated_at,
  COALESCE((
    SELECT array_agg(component_id) 
    FROM view_service_components 
    WHERE view_service_id = vs.id AND deleted_at IS NULL
  ), '{}') AS component_ids,
  COALESCE((
    SELECT array_agg(component_group_id) 
    FROM view_service_component_groups 
    WHERE view_service_id = vs.id AND deleted_at IS NULL
  ), '{}') AS component_group_ids
FROM
  view_services vs
WHERE
  vs.view_id = @view_id
  AND vs.service_id = @service_id
  AND vs.deleted_at IS NULL
LIMIT 1;
