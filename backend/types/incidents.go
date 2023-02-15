package types

import (
	"time"
)

type Incident struct {
	BaseModel
	Description        string
	Url                string
	IncidentCreatedAt  time.Time
	ProviderIncidentId string
	
}

type IncidentUpdate struct {
	BaseModel
	IncidentId  string
	Description string
	Status      string
	StatusTime  string
}

type IncidentUpdatesComponent struct {
	BaseModel
	IncidentUpdateId string
	ComponentId      string
}
