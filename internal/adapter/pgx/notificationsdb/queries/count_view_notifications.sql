SELECT COUNT(*)
FROM view_notifications
WHERE view_id = @view_id
  AND deleted_at IS NULL
  AND (
    @search::text = ''
    OR name ILIKE '%' || @search || '%'
    OR type ILIKE '%' || @search || '%'
  );
