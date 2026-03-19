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
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS incident_components;
-- +goose StatementEnd
