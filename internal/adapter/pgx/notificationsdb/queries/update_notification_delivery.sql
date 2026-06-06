UPDATE notification_deliveries
SET last_update_id = $1, updated_at = now()
WHERE id = $2;
