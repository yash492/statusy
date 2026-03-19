package schedulemaintenance

import (
	"time"

	"github.com/yash492/statusy/internal/common/nullable"
)

type ScheduleMaintenanceUpdate struct {
	ScheduleMaintenanceID         uint
	Description                   string
	ScheduleMaintenanceProviderID string
	ProviderID                    string
	Status                        string
	ProviderStatus                string
	StatusTime                    time.Time
}

type ScheduleMaintenance struct {
	Name              string
	Link              string
	ServiceID         uint
	StartsAt          time.Time
	EndsAt            time.Time
	ProviderImpact    nullable.Nullable[string]
	Impact            nullable.Nullable[string]
	ProviderID        string
	ProviderCreatedAt time.Time
	Updates           []ScheduleMaintenanceUpdate
}
