package incidents

import (
	"time"

	"github.com/yash492/statusy/internal/common"
	"github.com/yash492/statusy/internal/domain/components"
)

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
	Components        []components.Component
}
