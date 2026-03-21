WITH
    scheduled_maintenance_status_cte AS (
        SELECT
            scheduled_maintenances.id AS id,
            service_id AS service_id,
            scheduled_maintenances.title AS title,
            scheduled_maintenance_updates.status AS status,
            provider_created_at,
            link,
            starts_at,
            ends_at,
            RANK() OVER (
                PARTITION BY
                    scheduled_maintenance_updates.scheduled_maintenance_id
                ORDER BY scheduled_maintenance_updates.status_time DESC
            ) as rank_
        FROM scheduled_maintenances
            JOIN scheduled_maintenance_updates ON scheduled_maintenances.id = scheduled_maintenance_updates.scheduled_maintenance_id
        WHERE
            service_id = @service_id
    )
SELECT
    id,
    service_id,
    title,
    status,
    starts_at,
    ends_at
    link,
    provider_created_at
FROM scheduled_maintenance_status_cte
WHERE
    rank_ = 1
OFFSET
    @offset
LIMIT @limit