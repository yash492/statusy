package incidents

import (
	"time"

	"github.com/yash492/statusy/internal/common/nullable"
	"github.com/yash492/statusy/internal/domain/components"
)

type IncidentUpdate struct {
	IncidentID         uint
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
	ServiceID         uint
	ProviderImpact    nullable.Nullable[string]
	Impact            nullable.Nullable[string]
	ProviderID        string
	ProviderCreatedAt time.Time
	Updates           []IncidentUpdate
	Components        []components.Component
}
