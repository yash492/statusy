package schema

import (
	"database/sql"

	"github.com/google/uuid"
)

type ServicesForSubsciptions struct {
	ServiceID   uint   `db:"service_id"`
	ServiceName string `db:"service_name"`
}

type Subscription struct {
	BaseModel
	ServiceID       uint      `db:"service_id"`
	UUID            uuid.UUID `db:"uuid"`
	IsAllComponents bool      `db:"is_all_components"`
}

type SubscriptionWithComponents struct {
	ServiceID       uint          `db:"service_id"`
	ServiceName     string        `db:"service_name"`
	UUID            uuid.NullUUID `db:"uuid"`
	IsAllComponents sql.NullBool  `db:"is_all_components"`
	ComponentName   string        `db:"component_name"`
	ComponentID     uint          `db:"component_id"`
	IsConfigured    bool          `db:"is_configured"`
}
