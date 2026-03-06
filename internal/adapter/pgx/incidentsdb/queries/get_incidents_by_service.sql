SELECT 
    id,
    service_id,
    name AS title,
    COALESCE(impact, 'investigating') AS status,
    link AS url,
    created_at,
    updated_at
FROM incidents
WHERE service_id = @service_id
AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT @limit
OFFSET @offset