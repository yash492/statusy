INSERT INTO view_notifications (view_id, type, config, created_at, updated_at)
VALUES ($1, $2, $3, now(), now())
RETURNING id, view_id, type, config, created_at, updated_at;
