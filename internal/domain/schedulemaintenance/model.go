package schedulemaintenance

import "time"

type ScheduleMaintenance struct {
	Name       string
	ProviderID string
	StartsAt   time.Time
	EndsAt     time.Time
	Status     string
}
