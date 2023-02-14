package types

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Incident struct {
	gorm.Model
	Uuid               uuid.UUID
	Description        string
	Url                string
	IncidentCreatedAt  time.Time
	ProviderIncidentId string
}

type IncidentUpdate struct {
	gorm.Model
	Uuid        uuid.UUID
	IncidentId  string
	Description string
	Status      string
	StatusTime  string
}

type IncidentUpdatesComponent struct {
	gorm.Model
	IncidentUpdateId string
	ComponentId      string
}
