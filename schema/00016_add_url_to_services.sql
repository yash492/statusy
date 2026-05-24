-- +goose Up
-- +goose StatementBegin
ALTER TABLE services ADD COLUMN url TEXT NOT NULL DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE services DROP COLUMN url;
-- +goose StatementEnd
