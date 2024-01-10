package schema

import (
	"database/sql"

	"github.com/google/uuid"
)

type WebhookExtension struct {
	BaseModel
	UUID       uuid.UUID      `db:"uuid"`
	WebhookURL string         `db:"webhook_url"`
	Secret     sql.NullString `db:"secret"`
}
