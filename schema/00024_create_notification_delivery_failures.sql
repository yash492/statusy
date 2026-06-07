-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS notification_delivery_failures (
    id SERIAL PRIMARY KEY,
    view_notification_id INT NOT NULL REFERENCES view_notifications(id) ON DELETE CASCADE,
    alert_type TEXT NOT NULL, -- 'incident' or 'scheduled_maintenance'
    alert_id INT NOT NULL,
    update_id INT NOT NULL,
    error_message TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS notification_delivery_failures;
-- +goose StatementEnd
