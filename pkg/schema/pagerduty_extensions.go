package schema

import "github.com/google/uuid"

type PagerdutyExtension struct {
	BaseModel
	UUID       uuid.UUID `db:"uuid"`
	RoutingKey string    `db:"routing_key"`
}
