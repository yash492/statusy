-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX incident_updates_provider_id_idx ON incident_updates (provider_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS incident_updates_provider_id_idx;
-- +goose StatementEnd
