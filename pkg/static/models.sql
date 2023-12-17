CREATE TABLE IF NOT EXISTS services (
    id SERIAL NOT NULL,
    name TEXT NOT NULL,
    link TEXT NOT NULL,
    slug TEXT NOT NULL,
    provider_type TEXT NOT NULL,
    should_scrap_website boolean NOT NULL,
    incident_url TEXT NOT NULL,
    schedule_maintenance_url TEXT NOT NULL,
    components_url TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NOT NULL
  );

CREATE UNIQUE INDEX IF NOT EXISTS idx_services_slug ON services (slug); 

CREATE TABLE IF NOT EXISTS components (
    id SERIAL NOT NULL,
    name TEXT NOT NULL,
    service_id INTEGER NOT NULL,
    provider_component_id TEXT NOT NULL,
    created_at TIMESTAMPTZ NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_components_name_service_id ON components (service_id, name); 