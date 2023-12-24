package store

import "github.com/yash492/statusy/pkg/schema"

type ServiceStore interface {
	Create(services []schema.Service) error
	GetAll() ([]schema.Service, error)
}

type ComponentStore interface {
	Create(components []schema.Component) ([]schema.Component, error)
	GetAllByServiceID(serviceID uint) ([]schema.Component, error)
}

type IncidentStore interface {
	Create([]schema.Incident) ([]schema.Incident, error)
	GetByProviderID(providerID string) (schema.Incident, error)
	CreateIncidentUpdates(incidentUpdates []schema.IncidentUpdate) error
	GetLastIncidentUpdatesTimeByService(serviceID uint, incidentIDs []uint) ([]schema.LastIncidentUpdateForIncident, error)
	CreateIncidentComponents(incidentComponents []schema.IncidentComponent) error
}
