package models

import (
	"time"

	"gorm.io/gorm"
)

type Incident struct {
	BaseModel
	Description        string
	Url                string
	IncidentCreatedAt  time.Time
	ServiceId          uint
	ProviderIncidentId string
}

type IncidentUpdate struct {
	BaseModel
	IncidentId  uint
	Description string
	Status      string
	StatusTime  time.Time
}

// Using gorm.Model because incident_components table does not have uuid
type IncidentComponent struct {
	gorm.Model
	IncidentId  uint
	ComponentId uint
}

type LastUpdatedIncidentForSlug struct {
	Slug          string
	LastUpdatedAt time.Time
}
