CREATE TABLE IF NOT EXISTS services (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  slug TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX ON services (slug);

CREATE TABLE IF NOT EXISTS component_groups (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  service_id INT NOT NULL REFERENCES services (id),
  provider_id TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX ON component_groups (provider_id);

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

CREATE UNIQUE INDEX ON components (provider_id);

CREATE TABLE IF NOT EXISTS incidents (
  id SERIAL PRIMARY KEY,
  title TEXT NOT NULL,
  link TEXT NOT NULL,
  provider_impact TEXT,
  impact TEXT,
  service_id INT NOT NULL REFERENCES services (id),
  provider_id TEXT NOT NULL,
  provider_created_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOT NULL NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOT NULL NOW(),
  deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX ON incidents (provider_id);

CREATE TABLE IF NOT EXISTS incident_updates (
  id SERIAL PRIMARY KEY,
  incident_id INT NOT NULL REFERENCES incidents (id),
  description TEXT NOT NULL,
  provider_status TEXT NOT NULL,
  status TEXT NOT NULL,
  status_time TIMESTAMPTZ NOT NULL,
  provider_id TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX ON incident_updates (provider_id);

CREATE TABLE IF NOT EXISTS incident_components (
  id SERIAL PRIMARY KEY,
  incident_id INT NOT NULL REFERENCES incidents (id),
  component_id INT NOT NULL REFERENCES components (id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX ON incident_components (incident_id, component_id);

CREATE TABLE IF NOT EXISTS schedule_maintenances (
  id SERIAL PRIMARY KEY,
  title TEXT NOT NULL,
  link TEXT NOT NULL,
  service_id INT NOT NULL REFERENCES services (id),
  provider_id TEXT NOT NULL,
  scheduled_start_time TIMESTAMPTZ NOT NULL,
  scheduled_end_time TIMESTAMPTZ NOT NULL,
  provider_created_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);


 CREATE TABLE IF NOT EXISTS schedule_maintenance_updates (
    id SERIAL PRIMARY KEY,
    incident_id INT NOT NULL REFERENCES incidents (id),
    description TEXT NOT NULL,
    provider_status TEXT NOT NULL,
    status TEXT NOT NULL,
    status_time TIMESTAMPTZ NOT NULL,
    provider_id TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
