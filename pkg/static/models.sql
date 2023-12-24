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
    provider_id TEXT NOT NULL,
    created_at TIMESTAMPTZ NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_components_name_service_id ON components (service_id, name); 

CREATE TABLE IF NOT EXISTS incidents (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	link TEXT NOT NULL,
	service_id INT NOT NULL,
	provider_id TEXT NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	deleted_at TIMESTAMP WITH TIME ZONE,
	CONSTRAINT fk_service_id_incidents
		FOREIGN KEY(service_id) 
		REFERENCES services(id)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_provider_id ON incidents(provider_id);

CREATE TABLE IF NOT EXISTS incident_updates (
	id SERIAL PRIMARY KEY,
    incident_id INT NOT NULL,
	description TEXT NOT NULL,
	status TEXT NOT NULL,
	status_time TIMESTAMP WITH TIME ZONE NOT NULL,
	provider_id TEXT NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	deleted_at TIMESTAMP WITH TIME ZONE,
	CONSTRAINT fk_incident_id_incident_updates 
		FOREIGN KEY (incident_id)
		REFERENCES incidents(id)
);

CREATE TABLE IF NOT EXISTS incident_components (
	id SERIAL PRIMARY KEY,
  incident_id INT NOT NULL,
	component_id INT NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT fk_incident_components_incident_id
		FOREIGN KEY (incident_id)
		REFERENCES incidents(id),
    CONSTRAINT fk_incident_components_component_id
		FOREIGN KEY (component_id)
		REFERENCES components(id)
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_incident_id_component_id ON incident_components(incident_id, component_id);

