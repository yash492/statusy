UPDATE incidents i
SET is_resolved = COALESCE(
  (
    SELECT iu.status IN ('resolved', 'postmortem')
    FROM incident_updates iu
    WHERE iu.incident_id = i.id AND iu.deleted_at IS NULL
    ORDER BY iu.status_time DESC, iu.id DESC
    LIMIT 1
  ),
  FALSE
)
WHERE i.service_id = ANY(@service_ids) AND i.deleted_at IS NULL;
