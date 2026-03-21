package scheduledmaintenance

import (
	"time"

	"github.com/yash492/statusy/internal/common/nullable"
)

type ScheduledMaintenanceUpdate struct {
	scheduledMaintenanceID         uint
	Description                    string
	scheduledMaintenanceProviderID string
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
}
