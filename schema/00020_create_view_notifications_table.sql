-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS view_notifications (
    id SERIAL PRIMARY KEY,
    view_id INT NOT NULL REFERENCES views(id) ON DELETE CASCADE,
    type TEXT NOT NULL,
    config JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS view_notifications;
-- +goose StatementEnd
