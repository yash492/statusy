SELECT
    sm.id AS maintenance_id,
    smu.id AS update_id,
    sm.title,
    smu.status,
    smu.description,
    sm.provider_id,
    s.name AS service_name,
    COALESCE(
        json_agg(
            json_build_object(
                'name', c.name,
                'group_name', cg.name
            )
        ) FILTER (
            WHERE
                c.name IS NOT NULL
        ),
        '[]'
    ) AS components,
    sm.start_time,
    sm.end_time,
    smu.updated_at,
    sm.link
FROM
    scheduled_maintenance_updates smu
    JOIN scheduled_maintenances sm ON sm.id = smu.scheduled_maintenance_id
    JOIN services s ON s.id = sm.service_id
    LEFT JOIN scheduled_maintenance_components smc ON smc.scheduled_maintenance_id = sm.id
    LEFT JOIN components c ON c.id = smc.component_id
    LEFT JOIN component_groups cg ON cg.id = c.component_group_id
WHERE
    smu.id = $1
GROUP BY
    sm.id,
    smu.id,
    sm.title,
    smu.status,
    smu.description,
    sm.provider_id,
    s.name,
    sm.start_time,
    sm.end_time,
    smu.updated_at,
    sm.link;
