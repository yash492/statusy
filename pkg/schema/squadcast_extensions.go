package schema

import "github.com/google/uuid"

type SquadcastExtension struct {
	BaseModel
	UUID       uuid.UUID `db:"uuid"`
	WebhookURL string    `db:"webhook_url"`
}
