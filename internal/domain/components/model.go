package components

import "github.com/yash492/statusy/internal/domain/services"

type AggregateComponents struct {
	Service             services.Service
	GroupedComponents   []ComponentGroup
	UngroupedComponents []Component
}

type ComponentGroup struct {
	ID         uint
	Name       string
	ProviderID string
	Components []Component
}

type Component struct {
	ID         uint
	Name       string
	ProviderID string
}
