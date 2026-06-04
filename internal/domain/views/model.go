package views

import "time"

type View struct {
	ID          uint
	Name        string
	PublicID    string
	Description string
	IsDefault   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ViewServiceStatus struct {
	ID                           uint
	Name                         string
	Slug                         string
	Status                       string
	LastIncident                 string
	IncludeAllComponents         bool
	MonitorIncidents             bool
	MonitorScheduledMaintenances bool
	UpcomingMaintenance          string
	LastIncidentLink             string
	UpcomingMaintenanceLink      string
	ComponentIDs                 []int
	ComponentGroupIDs            []int
}

type ViewService struct {
	ID                           uint
	ViewID                       uint
	ServiceID                    uint
	IncludeAllComponents         bool
	MonitorIncidents             bool
	MonitorScheduledMaintenances bool
	ComponentIDs                 []int
	ComponentGroupIDs            []int
	CreatedAt                    time.Time
	UpdatedAt                    time.Time
}
