package store

import "github.com/yash492/statusy/pkg/schema"

type ServiceStore interface {
	Create(services []schema.Service) error
	GetAll() ([]schema.Service, error)
}

type ComponentStore interface {
	Create(components []schema.Component) error
}

type IncidentStore interface {
}
