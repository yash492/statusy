-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS scheduled_maintenances (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    link TEXT NOT NULL,
    provider_impact TEXT,
    impact TEXT,
    service_id INT NOT NULL REFERENCES services (id),
    provider_id TEXT NOT NULL,
    starts_at TIMESTAMPTZ NOT NULL,
    ends_at TIMESTAMPTZ NOT NULL,
    provider_created_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX scheduled_maintenance_provider_id_idx ON scheduled_maintenances (provider_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS scheduled_maintenances;
-- +goose StatementEnd