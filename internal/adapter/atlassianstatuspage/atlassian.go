package atlassianstatuspage

import "time"

type atlassianProvider struct {
	incidentUrl string
	serviceID   uint
}

type atlassianIncidentReq struct {
	Incidents []atlassianIncident `json:"incidents"`
}

type atlassianIncident struct {
	ID                 string                       `json:"id"`
	Name               string                       `json:"name"`
	Status             string                       `json:"status"`
	CreatedAt          time.Time                    `json:"created_at"`
	UpdatedAt          time.Time                    `json:"updated_at"`
	MonitoringAt       time.Time                    `json:"monitoring_at"`
	ResolvedAt         time.Time                    `json:"resolved_at"`
	Impact             string                       `json:"impact"`
	Shortlink          string                       `json:"shortlink"`
	StartedAt          time.Time                    `json:"started_at"`
	PageID             string                       `json:"page_id"`
	IncidentUpdates    []atlassianIncidentUpdate    `json:"incident_updates"`
	IncidentComponents []atlassianIncidentComponent `json:"components"`
}

type atlassianIncidentUpdate struct {
	ID         string    `json:"id"`
	Status     string    `json:"status"`
	Body       string    `json:"body"`
	IncidentID string    `json:"incident_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DisplayAt  time.Time `json:"display_at"`
}

type atlassianIncidentComponent struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type atlassianComponentsReq struct {
	Components []atlassianComponent `json:"components"`
}

type atlassianComponent struct {
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
