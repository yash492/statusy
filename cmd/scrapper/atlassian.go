package scrapper

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/jackc/pgx/v5"
	"github.com/samber/lo"
	"github.com/yash492/statusy/pkg/domain"
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

	_, err := client.R().SetHeader("Authorization", "OAuth 317d16908f244276a67169cc87dd65c0").SetResult(&atlassian).Get(a.incidentUrl)
	if err != nil {
		return err
	}

	//TODO: remove this
	// atlassian.Incidents = atlassian.Incidents[:4]

	components, err := domain.Component.GetAllByServiceID(a.serviceID)
	if err != nil {
		return err
	}

	componentsMap := lo.Associate(components, func(component schema.Component) (string, schema.Component) {
		return component.ProviderID, component
	})

	var newIncidents []schema.Incident
	var newIncidentsUpdates []schema.IncidentUpdate

	//This is a map between incident ID and incident update
	existingIncidentUpdateMap := make(map[uint][]schema.IncidentUpdate, 0)

	//This is a map between provider incident ID and incident update
	newIncidentUpdateMap := make(map[string][]atlassianIncidentUpdate, 0)

	//This a map between incident provider id and incident components
	newIncidentComponentMap := make(map[string][]atlassianIncidentComponent, 0)

	//This a map between incident id and incident components
	incidentComponentMap := make(map[uint][]atlassianIncidentComponent, 0)

	// This is start of incident parsing stage
	for _, incidentReq := range atlassian.Incidents {

		// Reversing the incident updates from atlassian as
		// it sends the payload in the following descending order resolved -> update -> investigating
		// which is not desirable since we need the updates in the ascending order
		// to perform connected dispatch events, for example squadcast can't resolve an incident
		// when an already resolved signal is sent before the trigger signal
		incidentUpdateReq := lo.Reverse(incidentReq.IncidentUpdates)

		// If incident already exists in the DB, fetch it from the DB
		// and the put it into the map for further scraping of incident updates and incident
		// components
		// This is to ensure bulk insert of incident updates
		existingIncident, err := domain.Incident.GetByProviderID(incidentReq.ID)
		if err == nil {
			incidentUpdates := a.parseExistingIncidentUpdates(incidentUpdateReq, existingIncident.ID)
			existingIncidentUpdateMap[existingIncident.ID] = incidentUpdates
			incidentComponentMap[existingIncident.ID] = incidentReq.IncidentComponents
			continue
		}

		newIncidents = append(newIncidents, schema.Incident{
			Name:           incidentReq.Name,
			Link:           incidentReq.Shortlink,
			ServiceID:      a.serviceID,
			ProviderImpact: incidentReq.Impact,
			Impact:         incidentReq.Impact,
			ProviderID:     incidentReq.ID,
			ProviderCreatedAt: incidentUpdateReq[0].CreatedAt,
		})

		// Using Provider incident ID since, a standart incident ID is not available
		newIncidentUpdateMap[incidentReq.ID] = incidentUpdateReq
		newIncidentComponentMap[incidentReq.ID] = incidentReq.IncidentComponents
	}

	insertedIncidents, err := domain.Incident.Create(newIncidents)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	//This adds incident ID to incident updates & incident components
	// from freshly inserted incidents
	for _, incident := range insertedIncidents {
		atlassianIncidentUpdates, ok := newIncidentUpdateMap[incident.ProviderID]
		// This should not happen
		if !ok {
			continue
		}

		atlassianIncidentComponent, ok := newIncidentComponentMap[incident.ProviderID]
		// This should not happen
		if !ok {
			continue
		}

		incidentComponentMap[incident.ID] = atlassianIncidentComponent

		incidentUpdated := lo.Map(atlassianIncidentUpdates, func(update atlassianIncidentUpdate, i int) schema.IncidentUpdate {

			status := normaliseProviderState(update.Status)

			// Assuming first incident update will be in a triggered state
			if i == 0 {
				if status != "resolved" && status != "postmortem" {
					status = types.IncidentTriggered
				}
			}

			return schema.IncidentUpdate{
				IncidentID:     incident.ID,
				Description:    update.Body,
				ProviderStatus: update.Status,
				StatusTime:     update.CreatedAt,
				ProviderID:     update.ID,
				Status:         status,
			}
		})

		newIncidentsUpdates = append(newIncidentsUpdates, incidentUpdated...)

	}

	newIncidentUpdatesFromExistingIncident, err := a.handleExistingIncidents(existingIncidentUpdateMap)
	if err != nil {
		return err
	}
	aggregatedIncidentUpdates := append(newIncidentUpdatesFromExistingIncident, newIncidentsUpdates...)
	insertedIncidentsUpdates, err := domain.Incident.CreateIncidentUpdates(aggregatedIncidentUpdates)
	if err != nil {
		return err
	}

	err = a.handleIncidentComponents(componentsMap, incidentComponentMap)
	if err != nil {
		return err
	}

	// Publish incident to dispatcher
	publishUpdatesToDispatcher(queue, insertedIncidentsUpdates)
	return nil
}

func (a atlassianProvider) handleExistingIncidents(existingIncidentUpdateMap map[uint][]schema.IncidentUpdate) ([]schema.IncidentUpdate, error) {
	incidentIDs := lo.Keys(existingIncidentUpdateMap)
	lastIncidentUpdates, err := domain.Incident.GetLastIncidentUpdatesTimeByService(a.serviceID, incidentIDs)
	if err != nil {
		return nil, err
	}

	newIncidentUpdates := make([]schema.IncidentUpdate, 0)

	for _, lastIncidentUpdate := range lastIncidentUpdates {
		if incidentUpdates, ok := existingIncidentUpdateMap[lastIncidentUpdate.IncidentID]; ok {
			for _, incidentUpdate := range incidentUpdates {
				if (!lastIncidentUpdate.LastIncidentUpdatesTime.Valid) || (incidentUpdate.StatusTime.After(lastIncidentUpdate.LastIncidentUpdatesTime.Time)) {
					newIncidentUpdates = append(newIncidentUpdates, schema.IncidentUpdate{
						IncidentID:     lastIncidentUpdate.IncidentID,
						Description:    incidentUpdate.Description,
						ProviderID:     incidentUpdate.ProviderID,
						ProviderStatus: incidentUpdate.ProviderStatus,
						StatusTime:     incidentUpdate.StatusTime,
						Status:         normaliseProviderState(incidentUpdate.ProviderStatus),
					})
				}
			}
		}
	}
	return newIncidentUpdates, nil
}

func (atlassianProvider) parseExistingIncidentUpdates(req []atlassianIncidentUpdate, incidentID uint) []schema.IncidentUpdate {
	incidentUpdates := lo.Map(req, func(req atlassianIncidentUpdate, _ int) schema.IncidentUpdate {
		incidentUpdate := schema.IncidentUpdate{
			Description:    req.Body,
			ProviderStatus: req.Status,
			StatusTime:     req.CreatedAt,
			IncidentID:     incidentID,
			ProviderID:     req.ID,
		}

		return incidentUpdate
	})

	return incidentUpdates

}

func (a atlassianProvider) handleIncidentComponents(componentsMap map[string]schema.Component, incidentComponentMap map[uint][]atlassianIncidentComponent) error {
	var incidentComponents []schema.IncidentComponent

	for incidentID, atlassianComponents := range incidentComponentMap {
		for _, atlassianComponent := range atlassianComponents {
			component, ok := componentsMap[atlassianComponent.ID]
			if !ok {
				returnedComponents, err := domain.Component.Create([]schema.Component{{
					Name:       atlassianComponent.Name,
					ProviderID: atlassianComponent.ID,
					ServiceID:  a.serviceID,
				}})

				if err != nil {
					return err
				}
				component = returnedComponents[0]
			}
			incidentComponents = append(incidentComponents, schema.IncidentComponent{
				IncidentID:  incidentID,
				ComponentID: component.ID,
			})
		}
	}

	return domain.Incident.CreateIncidentComponents(incidentComponents)
}

func publishUpdatesToDispatcher(dispatcherQueue *queue.Queue, incidentUpdates []schema.IncidentUpdate) {
	if dispatcherQueue == nil {
		return
	}
	for _, incidentUpdate := range incidentUpdates {
		dispatcherQueue.Publish(queue.IncidentPayload{
			State:          fmt.Sprintf("%v.%v", types.Incident, incidentUpdate.Status),
			IncidentUpdate: incidentUpdate,
		})
	}
}

func normaliseProviderState(providerState string) string {
	stateMap := map[string]string{
		"investigating": types.IncidentInProgress,
		"identified":    types.IncidentInProgress,
		"monitoring":    types.IncidentInProgress,
		"resolved":      types.IncidentResolved,
		"postmortem":    types.IncidentResolved,
	}
	return stateMap[providerState]
}
