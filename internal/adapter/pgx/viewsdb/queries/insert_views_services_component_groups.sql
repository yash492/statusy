INSERT INTO
  view_service_component_groups (
    view_service_id,
    component_group_id
  )
VALUES
  (
    @view_service_id,
    @component_group_id
  )
RETURNING
  *;
