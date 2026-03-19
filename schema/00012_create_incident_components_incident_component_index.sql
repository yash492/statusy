-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX incident_components_incident_component_idx ON incident_components (incident_id, component_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS incident_components_incident_component_idx;
-- +goose StatementEnd
