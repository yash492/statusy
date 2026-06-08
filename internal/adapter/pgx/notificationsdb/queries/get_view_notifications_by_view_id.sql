SELECT id, view_id, name, type, config, created_at, updated_at
FROM view_notifications
WHERE view_id = @view_id
  AND deleted_at IS NULL
  AND (
    @search::text = ''
    OR name ILIKE '%' || @search || '%'
    OR type ILIKE '%' || @search || '%'
  )
ORDER BY id DESC
LIMIT @limit
OFFSET @offset;
