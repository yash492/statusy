INSERT INTO
  view_service_components (
    view_service_id,
    component_id
  )
VALUES
  (
    @view_service_id,
    @component_id
  )
RETURNING
  *;
