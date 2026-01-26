-- name: CreateServices :copyfrom
INSERT INTO services (
    name,
    slug,
    provider_type,
    incidents_url,
    schedule_maintenances_url,
    components_url
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
);