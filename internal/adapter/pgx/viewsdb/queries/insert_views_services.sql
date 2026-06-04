INSERT INTO
  view_services (
    view_id,
    service_id,
    include_all_components,
    monitor_incidents,
    monitor_scheduled_maintenances
  )
VALUES
  (
    @view_id,
    @service_id,
    @include_all_components,
    @monitor_incidents,
    @monitor_scheduled_maintenances
  )
RETURNING
  *;