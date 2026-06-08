INSERT INTO view_notifications (view_id, name, type, config, created_at, updated_at)
VALUES ($1, $2, $3, $4, now(), now())
RETURNING id, view_id, name, type, config, created_at, updated_at;
