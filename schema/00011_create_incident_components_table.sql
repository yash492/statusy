-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS incident_components (
    id SERIAL PRIMARY KEY,
    incident_id INT NOT NULL REFERENCES incidents (id),
    component_id INT NOT NULL REFERENCES components (id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX incident_components_incident_component_idx ON incident_components (incident_id, component_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS incident_components;
-- +goose StatementEnd