-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS views (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    public_id TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS view_services (
    id SERIAL PRIMARY KEY,
    view_id INT NOT NULL REFERENCES views (id) ON DELETE CASCADE,
    service_id INT NOT NULL REFERENCES services (id) ON DELETE CASCADE,
    include_all_components BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);


CREATE UNIQUE INDEX view_services_view_service_partial_idx 
ON view_services (view_id, service_id) 
WHERE (deleted_at is NULL);

CREATE TABLE IF NOT EXISTS view_service_components (
    id SERIAL PRIMARY KEY,
    view_service_id INT NOT NULL REFERENCES view_services (id) ON DELETE CASCADE,
    component_id INT NOT NULL REFERENCES components (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS view_service_components_unique
ON view_service_components (view_service_id, component_id)
WHERE (deleted_at IS NULL);

CREATE TABLE IF NOT EXISTS view_service_component_groups (
    id SERIAL PRIMARY KEY,
    view_service_id INT NOT NULL REFERENCES view_services (id) ON DELETE CASCADE,
    component_group_id INT NOT NULL REFERENCES component_groups (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS view_service_component_groups_unique
ON view_service_component_groups (view_service_id, component_group_id)
WHERE (deleted_at IS NULL);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS view_service_component_groups;
DROP TABLE IF EXISTS view_service_components;
DROP TABLE IF EXISTS view_services;
DROP TABLE IF EXISTS views;
DROP INDEX IF EXISTS view_services_view_service_unique;
DROP INDEX IF EXISTS view_service_components_unique;
DROP INDEX IF EXISTS view_service_component_groups_unique;

-- +goose StatementEnd