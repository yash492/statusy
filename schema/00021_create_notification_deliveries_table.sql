-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS notification_deliveries (
    id SERIAL PRIMARY KEY,
    view_notification_id INT NOT NULL REFERENCES view_notifications(id) ON DELETE CASCADE,
    alert_type TEXT NOT NULL, -- 'incident' or 'scheduled_maintenance'
    alert_id INT NOT NULL,
    last_update_id INT NOT NULL,
    external_identifier TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT unique_notification_delivery UNIQUE (view_notification_id, alert_type, alert_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS notification_deliveries;
-- +goose StatementEnd
