package atlassian

import "time"

type IncidentReq struct {
	Incidents []Incident `json:"incidents"`
}

type ScheduledMaintenanceReq struct {
	ScheduledMaintenances []ScheduledMaintenance `json:"scheduled_maintenances"`
}

type Incident struct {
	ID                 string              `json:"id"`
	Name               string              `json:"name"`
	Status             string              `json:"status"`
	CreatedAt          time.Time           `json:"created_at"`
	UpdatedAt          time.Time           `json:"updated_at"`
	MonitoringAt       time.Time           `json:"monitoring_at"`
	ResolvedAt         time.Time           `json:"resolved_at"`
	Impact             string              `json:"impact"`
	Shortlink          string              `json:"shortlink"`
	StartedAt          time.Time           `json:"started_at"`
	PageID             string              `json:"page_id"`
	IncidentUpdates    []IncidentUpdate    `json:"incident_updates"`
	IncidentComponents []IncidentComponent `json:"components"`
}

type IncidentUpdate struct {
	ID         string    `json:"id"`
	Status     string    `json:"status"`
	Body       string    `json:"body"`
	IncidentID string    `json:"incident_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DisplayAt  time.Time `json:"display_at"`
}

type IncidentComponent struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type ComponentsReq struct {
	Components []Component `json:"components"`
}

type Component struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	Status             string    `json:"status"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	Position           int       `json:"position"`
	Description        *string   `json:"description"`
	Showcase           bool      `json:"showcase"`
	GroupID            *string   `json:"group_id"`
	PageID             string    `json:"page_id"`
	Group              bool      `json:"group"`
	OnlyShowIfDegraded bool      `json:"only_show_if_degraded"`
}

type ScheduledMaintenance struct {
	ID                 string              `json:"id"`
	Name               string              `json:"name"`
	Status             string              `json:"status"`
	CreatedAt          time.Time           `json:"created_at"`
	UpdatedAt          time.Time           `json:"updated_at"`
	MonitoringAt       time.Time           `json:"monitoring_at"`
	ResolvedAt         time.Time           `json:"resolved_at"`
	Impact             string              `json:"impact"`
	Shortlink          string              `json:"shortlink"`
	StartedAt          time.Time           `json:"started_at"`
	ScheduledFor       time.Time           `json:"scheduled_for"`
	ScheduledUntil     time.Time           `json:"scheduled_until"`
	PageID             string              `json:"page_id"`
	IncidentUpdates    []IncidentUpdate    `json:"incident_updates"`
	IncidentComponents []IncidentComponent `json:"components"`
}
