-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS scheduled_maintenance_components (
    id SERIAL PRIMARY KEY,
    scheduled_maintenance_id INT NOT NULL REFERENCES scheduled_maintenances (id),
    component_id INT NOT NULL REFERENCES components (id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX scheduled_maintenance_components_scheduled_maintenance_component_idx ON scheduled_maintenance_components (
    scheduled_maintenance_id,
    component_id
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS scheduled_maintenance_components;
-- +goose StatementEnd