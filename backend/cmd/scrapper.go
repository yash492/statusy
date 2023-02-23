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

	lastUpdatedIncident, err := incidentUpdatesEnv.Store.GetLastIncidentCreatedAtForSlug(a.ProviderSlug)
	if err != nil {
		return err
	}

	incidentUpdates := []models.IncidentUpdate{}
	incidentComponents := []models.IncidentComponent{}

	for _, incidentPayload := range a.Incidents.Incidents {

		if lastUpdatedIncident.LastUpdatedAt.After(incidentPayload.CreatedAt) {
			continue
		}

		incident, err := incidentsEnv.Store.AddIncident(models.Incident{
			Url:                incidentPayload.Shortlink,
			IncidentCreatedAt:  incidentPayload.CreatedAt,
			ProviderIncidentId: incidentPayload.ID,
			Description:        incidentPayload.Name,
			ServiceId:          service.ID,
		})
		if err != nil {
			return err
		}

		for _, updatePayload := range incidentPayload.IncidentUpdates {
			incidentUpdates = append(incidentUpdates, models.IncidentUpdate{
				IncidentId:  incident.ID,
				Description: updatePayload.Body,
				Status:      updatePayload.Status,
				StatusTime:  updatePayload.CreatedAt,
			})
		}

		// This get incident components, since the components are stored in incident updates
		// rather than the parent incident req payload
		incidentUpdateVar := incidentPayload.IncidentUpdates[0].AffectedComponents

		for _, componentPayload := range incidentUpdateVar {

			component, err := componentsEnv.Store.GetComponentsByCodeAndService(componentPayload.Code, service.ID)
			if err != nil {
				return err
			}

			incidentComponents = append(incidentComponents, models.IncidentComponent{
				IncidentId:  incident.ID,
				ComponentId: component.ID,
			})
		}
	}

	_, err = incidentUpdatesEnv.Store.AddIncidentUpdates(incidentUpdates)
	if err != nil {
		return err
	}
	_, err = incidentComponentsEnv.Store.AddIncidentComponents(incidentComponents)
	if err != nil {
		return err
	}

	return nil
}
