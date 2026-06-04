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
	ID                   uint
	Name                 string
	Slug                 string
	Status               string
	LastIncident         string
	IncludeAllComponents bool
}

type ViewService struct {
	ID                   uint
	ViewID               uint
	ServiceID            uint
	IncludeAllComponents bool
	ComponentIDs         []int
	ComponentGroupIDs    []int
	CreatedAt            time.Time
	UpdatedAt            time.Time
}
