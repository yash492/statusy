INSERT INTO notification_deliveries (view_notification_id, alert_type, alert_id, last_update_id, external_identifier, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, now(), now())
ON CONFLICT (view_notification_id, alert_type, alert_id) 
DO UPDATE SET last_update_id = EXCLUDED.last_update_id, external_identifier = EXCLUDED.external_identifier, updated_at = now();
