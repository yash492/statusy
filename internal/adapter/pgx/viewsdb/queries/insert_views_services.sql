INSERT INTO
  view_services (
    view_id,
    service_id,
    include_all_components
  )
VALUES
  (
    @view_id,
    @service_id,
    @include_all_components
  )
RETURNING
  *;