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
