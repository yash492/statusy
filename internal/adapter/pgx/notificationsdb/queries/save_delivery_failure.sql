INSERT INTO notification_delivery_failures (
    view_notification_id,
    alert_type,
    alert_id,
    update_id,
    error_message
) VALUES (
    $1, $2, $3, $4, $5
);
