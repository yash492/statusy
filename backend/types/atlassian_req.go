package types

import "time"

type AtlassianComponentsReq struct {
	Components []AtlassianComponent `json:"components"`
}

type AtlassianComponent struct {
	Name string `json:"name"`
}

type AtlassianStatusPageReq struct {
	Incidents []AtlassianIncident `json:"incidents"`
}

type AtlassianIncident struct {
	ID              string                    `json:"id"`
	Name            string                    `json:"name"`
	Status          string                    `json:"status"`
	CreatedAt       time.Time                 `json:"created_at"`
	UpdatedAt       time.Time                 `json:"updated_at"`
	MonitoringAt    time.Time                 `json:"monitoring_at"`
	ResolvedAt      time.Time                 `json:"resolved_at"`
	Impact          string                    `json:"impact"`
	Shortlink       string                    `json:"shortlink"`
	StartedAt       time.Time                 `json:"started_at"`
	PageID          string                    `json:"page_id"`
	IncidentUpdates []AtlassianIncidentUpdate `json:"incident_updates"`
}

type AtlassianIncidentUpdate struct {
	ID                 string                               `json:"id"`
	Status             string                               `json:"status"`
	Body               string                               `json:"body"`
	IncidentID         string                               `json:"incident_id"`
	CreatedAt          time.Time                            `json:"created_at"`
	UpdatedAt          time.Time                            `json:"updated_at"`
	DisplayAt          time.Time                            `json:"display_at"`
	AffectedComponents []AtlassianIncidentAffectedComponent `json:"affected_components"`
}

type AtlassianIncidentAffectedComponent struct {
	Code      string `json:"code"`
	Name      string `json:"name"`
	OldStatus string `json:"old_status"`
	NewStatus string `json:"new_status"`
}
