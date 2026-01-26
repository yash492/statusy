CREATE TABLE IF NOT EXISTS services (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    link TEXT NOT NULL,
    slug TEXT NOT NULL,
    provider_type TEXT NOT NULL,
    incident_url TEXT NOT NULL,
    schedule_maintenance_url TEXT NOT NULL,
    components_url TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
  );

CREATE TABLE IF NOT EXISTS component_groups (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    service_id INT NOT NULL,
    provider_id TEXT NOT NULL,
    created_at TIMESTAMPTZ NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT fk_service_id_component_groups
		FOREIGN KEY(service_id) 
		REFERENCES services(id)
);

CREATE TABLE IF NOT EXISTS components (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    service_id INT NOT NULL,
    provider_id TEXT NOT NULL,
    component_group_id INT,
    created_at TIMESTAMPTZ NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,    
	CONSTRAINT fk_service_id_components
		FOREIGN KEY(service_id) 
		REFERENCES services(id),
    CONSTRAINT fk_component_group_id_components
		FOREIGN KEY(component_group_id) 
		REFERENCES component_groups(id)
    
);


CREATE TABLE IF NOT EXISTS incidents (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	link TEXT NOT NULL,
	provider_impact TEXT,
	impact TEXT,
	service_id INT NOT NULL,
	provider_id TEXT NOT NULL,
	provider_created_at TIMESTAMPTZ NOT NULL,
	created_at TIMESTAMPTZ DEFAULT NOW(),
	updated_at TIMESTAMPTZ DEFAULT NOW(),
	deleted_at TIMESTAMPTZ,
	CONSTRAINT fk_service_id_incidents
		FOREIGN KEY(service_id) 
		REFERENCES services(id)
);

CREATE TABLE IF NOT EXISTS incident_updates (
	id SERIAL PRIMARY KEY,
    incident_id INT NOT NULL,
	description TEXT NOT NULL,
	provider_status TEXT NOT NULL,
	status TEXT NOT NULL,
	status_time TIMESTAMPTZ NOT NULL,
	provider_id TEXT NOT NULL,
	created_at TIMESTAMPTZ DEFAULT NOW(),
	updated_at TIMESTAMPTZ DEFAULT NOW(),
	deleted_at TIMESTAMPTZ,
	CONSTRAINT fk_incident_id_incident_updates 
		FOREIGN KEY (incident_id)
		REFERENCES incidents(id)
);

CREATE TABLE IF NOT EXISTS incident_components (
	id SERIAL PRIMARY KEY,
    incident_id INT NOT NULL,
	component_id INT NOT NULL,
	created_at TIMESTAMPTZ DEFAULT NOW(),
	updated_at TIMESTAMPTZ DEFAULT NOW(),
	deleted_at TIMESTAMPTZ,
    CONSTRAINT fk_incident_components_incident_id
		FOREIGN KEY (incident_id)
		REFERENCES incidents(id),
    CONSTRAINT fk_incident_components_component_id
		FOREIGN KEY (component_id)
		REFERENCES components(id)
);