package types

import (
	"time"
)

type Incident struct {
	BaseModel
	Description        string
	Url                string
	IncidentCreatedAt  time.Time
	ServiceId          int
	ProviderIncidentId string
}

type IncidentUpdate struct {
	BaseModel
	IncidentId  string
	Description string
	Status      string
	StatusTime  string
}

type IncidentsComponent struct {
	BaseModel
	IncidentId  string
	ComponentId string
}
