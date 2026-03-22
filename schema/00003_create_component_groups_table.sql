-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS component_groups (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    service_id INT NOT NULL REFERENCES services (id),
    provider_id TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS component_groups;
-- +goose StatementEnd