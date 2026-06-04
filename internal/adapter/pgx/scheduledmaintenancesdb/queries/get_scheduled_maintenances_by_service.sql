WITH
    scheduled_maintenance_status_cte AS (
        SELECT
            scheduled_maintenances.id AS id,
            scheduled_maintenances.service_id AS service_id,
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
            scheduled_maintenances.service_id = @service_id
            AND scheduled_maintenances.deleted_at IS NULL
            AND (
                @has_filter::boolean = false
                OR NOT EXISTS (SELECT 1 FROM scheduled_maintenance_components smc WHERE smc.scheduled_maintenance_id = scheduled_maintenances.id)
                OR EXISTS (
                    SELECT 1 FROM scheduled_maintenance_components smc
                    WHERE smc.scheduled_maintenance_id = scheduled_maintenances.id
                      AND (
                        smc.component_id = ANY(@component_ids::int[])
                        OR smc.component_id IN (
                          SELECT c.id FROM components c
                          WHERE c.component_group_id = ANY(@component_group_ids::int[])
                            AND c.deleted_at IS NULL
                        )
                      )
                )
            )
        ORDER BY provider_created_at DESC
    )
SELECT
    id,
    service_id,
    title,
    status,
    starts_at,
    ends_at,
    link,
    provider_created_at,
    COUNT(*) OVER() AS total_count
FROM scheduled_maintenance_status_cte
WHERE
    rank_ = 1
ORDER BY 
    provider_created_at DESC
OFFSET
    @offset
LIMIT @limit