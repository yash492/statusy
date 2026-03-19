-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS incident_updates (
  id SERIAL PRIMARY KEY,
  incident_id INT NOT NULL REFERENCES incidents (id),
  description TEXT NOT NULL,
  provider_status TEXT NOT NULL,
  status TEXT NOT NULL,
  status_time TIMESTAMPTZ NOT NULL,
  provider_id TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS incident_updates;
-- +goose StatementEnd
