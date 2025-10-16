package statuspage

import (
	"time"

	"github.com/yash492/statusy/internal/common"
)

type ProviderType string

type StatusPage interface {
	ScrapIncidents()
	// ScrapScheduleMaintainance()
	FetchComponents() (ComponentGroups, error)
	Slug() string
}

type ComponentGroups []ComponentGroups

type ComponentGroup struct {
	Name       string
	Components []Component
	// If there is no component grouping from the provider,
	// then it should show under any groups.
	IsGrouped bool
}

type Component struct {
	Name        string
	ServiceSlug string
	ProviderID  string
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
	Components        []Component
}
