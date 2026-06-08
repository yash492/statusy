SELECT COUNT(*)
FROM views
WHERE deleted_at IS NULL
  AND (
    @search::TEXT = ''
    OR name ILIKE '%' || @search || '%'
  );
