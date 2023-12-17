package store

import "github.com/yash492/statusy/pkg/schema"

type ServiceStore interface {
	Create(services []schema.Service) error
}

type IncidentStore interface {
}

type ComponentStore interface{}
