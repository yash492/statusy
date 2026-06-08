SELECT id, view_notification_id, alert_type, alert_id, last_update_id, external_identifier, created_at, updated_at
FROM notification_deliveries
WHERE view_notification_id = $1 AND alert_type = $2 AND alert_id = $3;
