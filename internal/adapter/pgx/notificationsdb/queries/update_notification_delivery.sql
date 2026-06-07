UPDATE notification_deliveries
SET last_update_id = $1, external_identifier = $2, updated_at = now()
WHERE id = $3;
