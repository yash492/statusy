WITH
    incident_status_cte AS (
        SELECT
            incidents.id AS id,
            incidents.service_id AS service_id,
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
            incidents.service_id = @service_id
            AND incidents.deleted_at IS NULL
            AND (
                @has_filter::boolean = false
                OR NOT EXISTS (SELECT 1 FROM incident_components ic WHERE ic.incident_id = incidents.id)
                OR EXISTS (
                    SELECT 1 FROM incident_components ic
                    WHERE ic.incident_id = incidents.id
                      AND (
                        ic.component_id = ANY(@component_ids::int[])
                        OR ic.component_id IN (
                          SELECT c.id FROM components c
                          WHERE c.component_group_id = ANY(@component_group_ids::int[])
                            AND c.deleted_at IS NULL
                        )
                      )
                )
            )
    )
SELECT
    id,
    service_id,
    title,
    status,
    link,
    provider_created_at,
    COUNT(*) OVER() AS total_count
FROM incident_status_cte
WHERE
    rank_ = 1
ORDER BY 
    provider_created_at DESC
OFFSET
    @offset
LIMIT @limit