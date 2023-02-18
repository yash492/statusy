package store

import (
	"backend/models"
)

type ServicesStore interface {
	AddServices([]models.Service) ([]models.Service, error)
	GetAllServices() ([]models.Service, error)
	GetServiceBySlug(slug string) (models.Service, error)
}

type ComponentsStore interface {
	AddComponents([]models.Component) ([]models.Component, error)
}

type IncidentStore interface {
	AddIncidents([]models.Incident) ([]models.Incident, error)
}

type IncidentUpdateStore interface {
	AddIncidentUpdates([]models.IncidentUpdate) ([]models.IncidentUpdate, error)
	GetLastIncidentCreatedAtForSlug(slug string) (models.LastUpdatedIncidentForSlug, error)
}

type IncidentComponentsStore interface {
	AddIncidentComponents([]models.IncidentsComponent) ([]models.IncidentsComponent, error)
}
