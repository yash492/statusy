package schema

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type ServicesForSubsciption struct {
	ServiceID   uint   `db:"service_id"`
	ServiceName string `db:"service_name"`
}

type Subscription struct {
	BaseModel
	ServiceID       uint      `db:"service_id"`
	UUID            uuid.UUID `db:"uuid"`
	IsAllComponents bool      `db:"is_all_components"`
}

type SubscriptionWithComponent struct {
	ServiceID     uint          `db:"service_id"`
	ServiceName   string        `db:"service_name"`
	UUID          uuid.NullUUID `db:"uuid"`
	ComponentName string        `db:"component_name"`
	ComponentID   uint          `db:"component_id"`
	IsConfigured  bool          `db:"is_configured"`
}

type SubscriptionForIncidentUpdate struct {
	ServiceID                    uint           `db:"service_id"`
	ServiceName                  string         `db:"service_name"`
	ComponentID                  uint           `db:"component_id"`
	ComponentName                string         `db:"component_name"`
	IncidentID                   uint           `db:"incident_id"`
	IncidentName                 string         `db:"incident_name"`
	IncidentLink                 string         `db:"incident_link"`
	IncidentImpact               sql.NullString `db:"incident_impact"`
	IncidentUpdate               string         `db:"incident_update"`
	IncidentUpdateID             uint           `db:"incident_update_id"`
	IncidentUpdateProviderStatus string         `db:"incident_update_provider_status"`
	IncidentUpdateStatus         string         `db:"incident_update_status"`
	IncidentUpdateStatusTime     time.Time      `db:"incident_update_status_time"`
	IsAllComponents              bool           `db:"is_all_components"`
}

type DashboardSubscription struct {
	SubscriptionsCount int64          `db:"subscriptions_count"`
	IncidentID         sql.NullInt64  `db:"incident_id"`
	ServiceID          uint           `db:"service_id"`
	ServiceName        string         `db:"service_name"`
	SubscriptionUUID   uuid.UUID      `db:"subscription_uuid"`
	IncidentName       sql.NullString `db:"incident_name"`
	IncidentLink       sql.NullString `db:"incident_link"`
	IncidentImpact     sql.NullString `db:"incident_impact"`
	IsDown             bool           `db:"is_down"`
}

type SubscriptionIncident struct {
	TotalCount                int64          `db:"total_count"`
	IncidentID                sql.NullInt64  `db:"incident_id"`
	LastUpdatedStatusTime     sql.NullTime   `db:"last_updated_status_time"`
	IncidentStatus            sql.NullString `db:"incident_status"`
	IncidentNormalisedStatus  sql.NullString `db:"incident_normalised_status"`
	IncidentCreatedAt         sql.NullTime   `db:"incident_created_at"`
	IncidentName              sql.NullString `db:"incident_name"`
	IncidentLink              sql.NullString `db:"incident_link"`
	ServiceID                 uint           `db:"service_id"`
	ServiceName               string         `db:"service_name"`
	IsAllComponentsConfigured bool           `db:"is_all_components_configured"`
}
