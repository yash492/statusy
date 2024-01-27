package schema

import (
	"database/sql"
	"time"
)

type Incident struct {
	//
	BaseModel
	Name              string         `db:"name"`
	Link              string         `db:"link"`
	ProviderImpact    sql.NullString `db:"provider_impact"`
	Impact            sql.NullString `db:"impact"`
	ServiceID         uint           `db:"service_id"`
	ProviderID        string         `db:"provider_id"`
	ProviderCreatedAt time.Time      `db:"provider_created_at"`
}

type IncidentUpdate struct {
	BaseModel
	IncidentID     uint      `db:"incident_id"`
	Description    string    `db:"description"`
	ProviderID     string    `db:"provider_id"`
	Status         string    `db:"status"`
	ProviderStatus string    `db:"provider_status"`
	StatusTime     time.Time `db:"status_time"`
}

type IncidentComponent struct {
	BaseModel
	IncidentID  uint `db:"incident_id"`
	ComponentID uint `db:"component_id"`
}

type LastIncidentUpdateForIncident struct {
	LastIncidentUpdatesTime sql.NullTime `db:"last_incident_updates_time"`
	IncidentID              uint         `db:"incident_id"`
}
