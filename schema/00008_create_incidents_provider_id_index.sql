-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX incidents_provider_id_idx ON incidents (provider_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS incidents_provider_id_idx;
-- +goose StatementEnd
