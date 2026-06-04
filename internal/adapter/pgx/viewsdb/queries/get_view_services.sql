SELECT
  s.id,
  s.name,
  s.slug,
  vs.include_all_components
FROM
  view_services vs
  JOIN services s ON s.id = vs.service_id
WHERE
  vs.view_id = @view_id
  AND vs.deleted_at IS NULL
ORDER BY
  vs.id ASC;
