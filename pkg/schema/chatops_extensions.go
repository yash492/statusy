package schema

import "github.com/google/uuid"

type ChatopsExtension struct {
	BaseModel
	UUID       uuid.UUID `db:"uuid"`
	Type       string    `db:"type"`
	WebhookURL string    `db:"webhook_url" json:"webhook_url"`
}
