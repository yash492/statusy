SELECT id, name, slug, incidents_url, schedule_maintenances_url, components_url, provider_type, created_at, updated_at, deleted_at
FROM services
WHERE slug = @slug AND deleted_at IS NULL