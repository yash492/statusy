-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS incidents (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    link TEXT NOT NULL,
    provider_impact TEXT,
    impact TEXT,
    service_id INT NOT NULL REFERENCES services (id),
    provider_id TEXT NOT NULL,
    provider_created_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS incidents;
-- +goose StatementEnd