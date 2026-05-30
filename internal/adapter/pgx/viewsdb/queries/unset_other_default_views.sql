UPDATE views
SET is_default = FALSE, updated_at = now()
WHERE id <> @id AND deleted_at IS NULL;
