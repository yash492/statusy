package cmd

import (
	"backend/models"
	"backend/types"
	"encoding/json"
	"io"
	"net/http"
)

type Scrapper interface {
	ScrapIncidents() error
	ScrapMaintenance() error
}

type AtlassianIncidents struct {
	IncidentUrl  string
	Incidents    types.AtlassianStatusPageReq
	ProviderSlug string
}

func (a *AtlassianIncidents) ScrapIncidents() error {

	resp, err := http.Get(a.IncidentUrl)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &a.Incidents)
	if err != nil {
		return err
	}

	service, err := servicesEnv.Store.GetServiceBySlug(a.ProviderSlug)
	if err != nil {
		return err
	}

	incidents := []models.Incident{}
	lastUpdatedIncident, err := incidentUpdatesEnv.Store.GetLastIncidentCreatedAtForSlug(a.ProviderSlug)
	if err != nil {
		return err
	}

	for _, incident := range a.Incidents.Incidents {

		if lastUpdatedIncident.LastUpdatedAt.After(incident.CreatedAt) {
			continue
		}

		incidents = append(incidents, models.Incident{
			Url:                incident.Shortlink,
			IncidentCreatedAt:  incident.CreatedAt,
			ProviderIncidentId: incident.ID,
			Description:        incident.Name,
			ServiceId:          service.ID,
		})
	}

	_, err = incidentsEnv.Store.AddIncidents(incidents)
	if err != nil {
		return err
	}

	return nil
}
