package scheduledmaintenance

import (
	"time"

	"github.com/yash492/statusy/internal/common/nullable"
	"github.com/yash492/statusy/internal/domain/components"
)

type ScheduledMaintenanceUpdate struct {
	ScheduledMaintenanceID         uint
	Description                    string
	ScheduledMaintenanceProviderID string
	ProviderID                     string
	Status                         string
	ProviderStatus                 string
	StatusTime                     time.Time
}

type ScheduledMaintenance struct {
	Name              string
	Link              string
	ServiceID         uint
	StartsAt          time.Time
	EndsAt            time.Time
	ProviderImpact    nullable.Nullable[string]
	Impact            nullable.Nullable[string]
	ProviderID        string
	ProviderCreatedAt time.Time
	Updates           []ScheduledMaintenanceUpdate
	Components        []components.Component
}
