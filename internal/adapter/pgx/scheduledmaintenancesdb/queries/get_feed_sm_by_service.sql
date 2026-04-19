WITH
    sm_status_cte AS (
        SELECT
            scheduled_maintenances.id AS id,
            service_id AS service_id,
            scheduled_maintenances.title AS title,
            scheduled_maintenance_updates.status AS status,
            provider_created_at,
            link,
            RANK() OVER (
                PARTITION BY scheduled_maintenance_updates.scheduled_maintenance_id
                ORDER BY scheduled_maintenance_updates.status_time DESC
            ) as rank_
        FROM scheduled_maintenances
            JOIN scheduled_maintenance_updates ON scheduled_maintenances.id = scheduled_maintenance_updates.scheduled_maintenance_id
        WHERE
            service_id = @service_id AND scheduled_maintenances.deleted_at IS NULL
    )
SELECT
    id,
    service_id,
    title,
    status,
    link,
    provider_created_at,
    COALESCE(
        (
            SELECT string_agg(c.name, ', ')
            FROM scheduled_maintenance_components smc
            JOIN components c ON smc.component_id = c.id
            WHERE smc.scheduled_maintenance_id = sm_status_cte.id
        ),
        ''
    ) AS affected_components
FROM sm_status_cte
WHERE rank_ = 1
ORDER BY provider_created_at DESC
LIMIT @limit OFFSET @offset
