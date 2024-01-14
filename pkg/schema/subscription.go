package schema

import (
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
	ServiceID     uint          `db:"service_id"`
	ServiceName   string        `db:"service_name"`
	UUID          uuid.NullUUID `db:"uuid"`
	ComponentName string        `db:"component_name"`
	ComponentID   uint          `db:"component_id"`
	IsConfigured  bool          `db:"is_configured"`
}

type SubscriptionForIncidentUpdates struct {
	ServiceID                    uint   `db:"service_id"`
	ServiceName                  string `db:"service_name"`
	ComponentID                  uint   `db:"component_id"`
	ComponentName                string `db:"component_name"`
	IncidentID                   uint   `db:"incident_id"`
	IncidentName                 string `db:"incident_name"`
	IncidentLink                 string `db:"incident_link"`
	IncidentImpact               string `db:"incident_impact"`
	IncidentUpdate               string `db:"incident_update"`
	IncidentUpdateProviderStatus string `db:"incident_update_provider_status"`
	IncidentUpdateStatus         string `db:"incident_update_status"`
	IsAllComponents              bool   `db:"is_all_components"`
}
