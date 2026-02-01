package statuspage

import (
	"time"

	"github.com/yash492/statusy/internal/common"
)

type ProviderType string
type ServiceSlug string

type StatusPageProvider interface {
	ScrapIncidents() ([]Incident, error)
	// ScrapScheduleMaintainance()
	ScrapComponents() (AggregateComponents, error)
	Slug() ServiceSlug
}

type Service struct {
	ID   uint
	Name string
}

type AggregateComponents struct {
	Service             Service
	GroupedComponents   []ComponentGroup
	UngroupedComponents []Component
}

type ComponentGroup struct {
	Name       string
	ProviderID string
	Components []Component
}

type Component struct {
	Name       string
	ProviderID string
}

type IncidentUpdate struct {
	Description        string
	IncidentProviderID string
	ProviderID         string
	Status             string
	ProviderStatus     string
	StatusTime         time.Time
}

type Incident struct {
	Name              string
	Link              string
	ServiceSlug       string
	ProviderImpact    common.Nullable[string]
	Impact            common.Nullable[string]
	ProviderID        string
	ProviderCreatedAt time.Time
	Updates           []IncidentUpdate
	Components        []Component
}

func (s ServiceSlug) String() string {
	return string(s)
}

func (s ProviderType) String() string {
	return string(s)
}
