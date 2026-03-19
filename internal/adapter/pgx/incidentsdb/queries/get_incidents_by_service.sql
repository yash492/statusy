WITH
    incident_status_cte AS (
        SELECT
            incidents.id AS id,
            service_id AS service_id,
            incidents.title AS title,
            incident_updates.status AS status,
            provider_created_at,
            link,
            RANK() OVER (
                PARTITION BY
                    incident_updates.incident_id
                ORDER BY incident_updates.status_time DESC
            ) as rank_
        FROM incidents
            JOIN incident_updates ON incidents.id = incident_updates.incident_id
        WHERE
            service_id = @service_id
    )
SELECT
    id,
    service_id,
    title,
    status,
    link,
    provider_created_at
FROM incident_status_cte
WHERE
    rank_ = 1
OFFSET
    @offset
LIMIT @limit