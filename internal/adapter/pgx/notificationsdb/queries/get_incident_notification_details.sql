SELECT
    i.id AS incident_id,
    iu.id AS update_id,
    i.title,
    iu.status,
    iu.description,
    i.provider_id,
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
    iu.updated_at,
    i.link
FROM
    incident_updates iu
    JOIN incidents i ON i.id = iu.incident_id
    JOIN services s ON s.id = i.service_id
    LEFT JOIN incident_components ic ON ic.incident_id = i.id
    LEFT JOIN components c ON c.id = ic.component_id
    LEFT JOIN component_groups cg ON cg.id = c.component_group_id
WHERE
    iu.id = $1
GROUP BY
    i.id,
    iu.id,
    i.title,
    iu.status,
    iu.description,
    i.provider_id,
    s.name,
    iu.updated_at,
    i.link;
