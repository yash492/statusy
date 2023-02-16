package store

import (
	"backend/types"
)

type ServicesStore interface {
	AddServices([]types.Service) ([]types.Service, error)
	GetAllServices() ([]types.Service, error)
}

type ComponentsStore interface {
	AddComponents([]types.Component) ([]types.Component, error)
}

type IncidentStore interface {
	AddIncidents([]types.Incident) ([]types.Incident, error)
}

type IncidentUpdateStore interface {
	AddIncidentUpdates([]types.IncidentUpdate) ([]types.IncidentUpdate, error)
}

type IncidentUpdateComponentsStore interface {
	AddIncidentUpdatesComponents([]types.IncidentUpdatesComponent) ([]types.IncidentUpdatesComponent, error)
}
