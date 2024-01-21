CREATE TABLE IF NOT EXISTS services (
    id SERIAL PRIMARY KEY,
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
    deleted_at TIMESTAMPTZ
  );

CREATE UNIQUE INDEX IF NOT EXISTS idx_services_slug ON services (slug); 

CREATE TABLE IF NOT EXISTS components (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    service_id INT NOT NULL,
    provider_id TEXT NOT NULL,
    created_at TIMESTAMPTZ NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
	CONSTRAINT fk_service_id_components
		FOREIGN KEY(service_id) 
		REFERENCES services(id)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_components_name_service_id ON components (service_id, name); 

CREATE TABLE IF NOT EXISTS incidents (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	link TEXT NOT NULL,
	provider_impact TEXT NOT NULL,
	impact TEXT NOT NULL,
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

CREATE UNIQUE INDEX IF NOT EXISTS idx_provider_id ON incidents(provider_id);

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
CREATE UNIQUE INDEX IF NOT EXISTS idx_incident_id_component_id ON incident_components(incident_id, component_id);

CREATE TABLE IF NOT EXISTS subscriptions (
	id SERIAL PRIMARY KEY,
	uuid UUID DEFAULT (uuid_generate_v4()),
	service_id INT NOT NULL,
	is_all_components BOOLEAN NOT NULL,
	created_at TIMESTAMPTZ DEFAULT NOW(),
	updated_at TIMESTAMPTZ DEFAULT NOW(),
	deleted_at TIMESTAMPTZ,
	CONSTRAINT fk_subscriptions_service_id
		FOREIGN KEY (service_id)
		REFERENCES services(id)
		
);

CREATE TABLE IF NOT EXISTS subscription_components (
	id SERIAL PRIMARY KEY,
	subscription_id INT NOT NULL,
	component_id INT NOT NULL,
	created_at TIMESTAMPTZ DEFAULT NOW(),
	updated_at TIMESTAMPTZ DEFAULT NOW(),
	deleted_at TIMESTAMPTZ,
	CONSTRAINT fk_subscriptions_component_id
		FOREIGN KEY (component_id)
		REFERENCES components(id),
	CONSTRAINT fk_subscription_components_subscription_id
		FOREIGN KEY (subscription_id)
		REFERENCES subscriptions(id)
		ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS squadcast_extensions (
	id SERIAL PRIMARY KEY,
	uuid UUID NOT NULL DEFAULT (uuid_generate_v4()),
	webhook_url TEXT NOT NULL,
	created_at TIMESTAMPTZ DEFAULT NOW(),
	updated_at TIMESTAMPTZ DEFAULT NOW(),
	deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_squadcast_extension_uuid ON squadcast_extensions(uuid); 


CREATE TABLE IF NOT EXISTS pagerduty_extensions (
	id SERIAL PRIMARY KEY,
	uuid UUID NOT NULL DEFAULT (uuid_generate_v4()),
	routing_key TEXT NOT NULL,
	created_at TIMESTAMPTZ DEFAULT NOW(),
	updated_at TIMESTAMPTZ DEFAULT NOW(),
	deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_pagerduty_extension_uuid ON pagerduty_extensions(uuid); 


CREATE TABLE IF NOT EXISTS chatops_extensions (
	id SERIAL PRIMARY KEY,
	uuid UUID NOT NULL DEFAULT (uuid_generate_v4()),
	-- type - slack, msteams, discord
	type TEXT NOT NULL,
	webhook_url TEXT NOT NULL,
	created_at TIMESTAMPTZ DEFAULT NOW(),
	updated_at TIMESTAMPTZ DEFAULT NOW(),
	deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_chatops_extension_uuid ON chatops_extensions(uuid); 

CREATE TABLE IF NOT EXISTS webhook_extensions (
	id SERIAL PRIMARY KEY,
	uuid UUID NOT NULL DEFAULT (uuid_generate_v4()),
	webhook_url TEXT NOT NULL,
	secret TEXT,
	created_at TIMESTAMPTZ DEFAULT NOW(),
	updated_at TIMESTAMPTZ DEFAULT NOW(),
	deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_webhooks_extension_uuid ON webhook_extensions(uuid); 