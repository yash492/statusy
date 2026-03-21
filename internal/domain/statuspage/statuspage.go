package statuspage

import (
	"github.com/yash492/statusy/internal/domain/components"
	"github.com/yash492/statusy/internal/domain/incidents"
	"github.com/yash492/statusy/internal/domain/services"
)

type StatusPageProvider interface {
	ScrapIncidents() ([]incidents.Incident, error)
	// ScrapscheduledMaintenance()
	ScrapComponents() (components.AggregateComponents, error)
	Slug() services.ServiceSlug
	ID() uint
	Name() string
	NewWithServiceID(id uint) StatusPageProvider
}
