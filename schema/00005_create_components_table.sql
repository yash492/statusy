-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS components (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  service_id INT NOT NULL,
  provider_id TEXT NOT NULL,
  component_group_id INT REFERENCES component_groups (id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS components;
-- +goose StatementEnd
