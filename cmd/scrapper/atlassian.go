package scrapper

import (
	"database/sql"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/samber/lo"
	"github.com/yash492/statusy/pkg/queue"
	"github.com/yash492/statusy/pkg/schema"
	"github.com/yash492/statusy/pkg/types"
)

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

func (a atlassianProvider) scrap(client *resty.Client, queue *queue.Queue) error {
	var atlassian atlassianIncidentReq

	_, err := client.R().SetResult(&atlassian).Get(a.incidentUrl)
	if err != nil {
		return err
	}

	statusPageIncidents := a.normaliseIncidents(atlassian.Incidents)
	return storeAndDispatchIncidents(queue, statusPageIncidents)

}

func (a atlassianProvider) normaliseIncidents(incidentsReq []atlassianIncident) []StatusPageIncident {
	statusPageIncidents := make([]StatusPageIncident, 0)

	for _, incidentReq := range incidentsReq {
		impact := sql.NullString{
			String: incidentReq.Impact,
			Valid:  true,
		}

		incidentUpdates := make([]schema.IncidentUpdate, 0)

		incident := schema.Incident{
			Name:              incidentReq.Name,
			Link:              incidentReq.Shortlink,
			ProviderImpact:    impact,
			Impact:            impact,
			ServiceID:         a.serviceID,
			ProviderID:        incidentReq.ID,
			ProviderCreatedAt: incidentReq.CreatedAt,
		}

		// Reversing the incident updates from atlassian as
		// it sends the payload in the following descending order resolved -> update -> investigating
		// which is not desirable since we need the updates in the ascending order
		// to perform connected dispatch events, for example squadcast can't resolve an incident
		// when an already resolved signal is sent before the trigger signal
		incidentReqUpdates := lo.Reverse(incidentReq.IncidentUpdates)
		for i, incidentUpdateReq := range incidentReqUpdates {

			status := a.normaliseProviderState(incidentUpdateReq.Status)
			if i == 0 && incidentUpdateReq.Status != "resolved" && incidentUpdateReq.Status != "postmortem" {
				status = types.IncidentTriggered
			}

			incidentUpdate := schema.IncidentUpdate{
				Description:    incidentUpdateReq.Body,
				ProviderID:     incidentUpdateReq.ID,
				Status:         status,
				ProviderStatus: incidentUpdateReq.Status,
				StatusTime:     incidentUpdateReq.CreatedAt,
			}
			incidentUpdates = append(incidentUpdates, incidentUpdate)

		}

		incidentComponents := lo.Map(incidentReq.IncidentComponents, func(incidentComponent atlassianIncidentComponent, _ int) StatusPageIncidentComponent {
			return StatusPageIncidentComponent{
				ProviderComponentID: incidentComponent.ID,
				ComponentName:       incidentComponent.Name,
			}
		})

		statusPageIncident := StatusPageIncident{
			Incident:           incident,
			IncidentUpdates:    incidentUpdates,
			IncidentComponents: incidentComponents,
		}
		statusPageIncidents = append(statusPageIncidents, statusPageIncident)
	}

	return statusPageIncidents
}

func (atlassianProvider) normaliseProviderState(providerState string) string {
	stateMap := map[string]string{
		"investigating": types.IncidentInProgress,
		"identified":    types.IncidentInProgress,
		"monitoring":    types.IncidentInProgress,
		"resolved":      types.IncidentResolved,
		"postmortem":    types.IncidentResolved,
	}
	return stateMap[providerState]
}
