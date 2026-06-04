-- +goose Up
-- +goose StatementBegin
ALTER TABLE incidents ADD COLUMN is_resolved BOOLEAN NOT NULL DEFAULT FALSE;

-- Index only active/unresolved incidents (typically 0 to 5 rows in the database)
CREATE INDEX IF NOT EXISTS incidents_active_idx ON incidents (service_id) WHERE NOT is_resolved;

-- Index upcoming/ongoing scheduled maintenances
CREATE INDEX IF NOT EXISTS scheduled_maintenances_active_idx ON scheduled_maintenances (service_id, ends_at DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS scheduled_maintenances_active_idx;
DROP INDEX IF EXISTS incidents_active_idx;
ALTER TABLE incidents DROP COLUMN is_resolved;
-- +goose StatementEnd
