-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX component_groups_provider_id_idx ON component_groups (provider_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS component_groups_provider_id_idx;
-- +goose StatementEnd