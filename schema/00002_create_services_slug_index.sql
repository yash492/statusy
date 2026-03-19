-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX services_slug_idx ON services (slug);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS services_slug_idx;
-- +goose StatementEnd
