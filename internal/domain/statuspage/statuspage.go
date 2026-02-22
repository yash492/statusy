package statuspage

import (
	"github.com/yash492/statusy/internal/domain/components"
	"github.com/yash492/statusy/internal/domain/incidents"
	"github.com/yash492/statusy/internal/domain/services"
)

type ProviderType = services.ProviderType
type ServiceSlug = services.ServiceSlug

type StatusPageProvider interface {
	ScrapIncidents() ([]incidents.Incident, error)
	// ScrapScheduleMaintainance()
	ScrapComponents() (components.AggregateComponents, error)
	Slug() ServiceSlug
	ID() uint
	Name() string
}
