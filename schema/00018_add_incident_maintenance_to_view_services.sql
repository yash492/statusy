-- +goose Up
-- +goose StatementBegin
ALTER TABLE view_services
ADD COLUMN monitor_incidents BOOLEAN NOT NULL DEFAULT TRUE,
ADD COLUMN monitor_scheduled_maintenances BOOLEAN NOT NULL DEFAULT TRUE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE view_services
DROP COLUMN monitor_incidents,
DROP COLUMN monitor_scheduled_maintenances;
-- +goose StatementEnd
