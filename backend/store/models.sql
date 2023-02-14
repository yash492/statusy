CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- services
CREATE TABLE IF NOT EXISTS services (
	id SERIAL PRIMARY KEY,
	uuid UUID DEFAULT (uuid_generate_v4()),
	name TEXT NOT NULL,
	link TEXT NOT NULL,
	slug TEXT NOT NULL,
	provider_type TEXT NOT NULL,
	should_scrap_website BOOLEAN NOT NULL,
	incident_url TEXT NOT NULL,
	schedule_maintenance_url TEXT NOT NULL,
	components_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX IF NOT EXISTS unique_slug_index ON services (slug); 

CREATE TABLE IF NOT EXISTS components (
	id SERIAL PRIMARY KEY,
	uuid UUID DEFAULT (uuid_generate_v4()),
  	name TEXT NOT NULL,
	slug TEXT NOT NULL,
	service_id INT NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	deleted_at TIMESTAMP WITH TIME ZONE,
	CONSTRAINT fk_service_id_components 
		FOREIGN KEY(service_id) 
		REFERENCES services(id)
);

CREATE UNIQUE INDEX IF NOT EXISTS unique_slug_index_components ON components (service_id, slug); 


CREATE TABLE IF NOT EXISTS incidents (
	id SERIAL PRIMARY KEY,
	uuid UUID DEFAULT (uuid_generate_v4()),
	description TEXT NOT NULL,
	url TEXT NOT NULL,
	incident_created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    provider_incident_id TEXT NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS incident_updates (
	id SERIAL PRIMARY KEY,
	uuid UUID DEFAULT (uuid_generate_v4()),
    incident_id INT NOT NULL,
	description TEXT NOT NULL,
	status TEXT NOT NULL,
	status_time TIMESTAMP WITH TIME ZONE NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	deleted_at TIMESTAMP WITH TIME ZONE,
	CONSTRAINT fk_incident_id_incident_updates 
		FOREIGN KEY (incident_id)
		REFERENCES incidents(id)
);

CREATE TABLE IF NOT EXISTS incident_updates_components (
	id SERIAL PRIMARY KEY,
    incident_update_id INT NOT NULL,
	component_id INT NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT fk_incident_updates_components_incident_update_id
		FOREIGN KEY (incident_update_id)
		REFERENCES incident_updates(id),
    CONSTRAINT fk_incident_updates_components_component_id
		FOREIGN KEY (component_id)
		REFERENCES components(id)
);