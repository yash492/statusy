-- +goose Up
-- +goose StatementBegin
ALTER TABLE view_notifications ADD COLUMN name TEXT NOT NULL DEFAULT '';
ALTER TABLE view_notifications ADD CONSTRAINT uq_view_notifications_view_id_name UNIQUE (view_id, name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE view_notifications DROP CONSTRAINT uq_view_notifications_view_id_name;
ALTER TABLE view_notifications DROP COLUMN name;
-- +goose StatementEnd
