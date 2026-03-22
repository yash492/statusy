-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX components_provider_id_idx ON components (provider_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS components_provider_id_idx;
-- +goose StatementEnd