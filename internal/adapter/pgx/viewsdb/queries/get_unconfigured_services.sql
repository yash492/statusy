SELECT
  id,
  name,
  slug,
  url
FROM
  services
WHERE
  deleted_at IS NULL
  AND id NOT IN (
    SELECT service_id
    FROM view_services
    WHERE view_id = @view_id AND deleted_at IS NULL
  )
  AND (
    @search::TEXT = ''
    OR name ILIKE '%' || @search || '%'
    OR slug ILIKE '%' || @search || '%'
  )
ORDER BY
  name ASC
LIMIT 5;

