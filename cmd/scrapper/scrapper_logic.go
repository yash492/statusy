package scrapper

import (
	"fmt"

	"github.com/samber/lo"
	"github.com/yash492/statusy/pkg/domain"
	"github.com/yash492/statusy/pkg/queue"
	"github.com/yash492/statusy/pkg/schema"
	"github.com/yash492/statusy/pkg/types"
)

type StatusPageIncident struct {
	Incident           schema.Incident
	IncidentUpdates    []schema.IncidentUpdate
	IncidentComponents []StatusPageIncidentComponent
}

type StatusPageIncidentComponent struct {
	IncidentComponents  schema.IncidentComponent
	ProviderComponentID string
	ComponentName       string
}

func storeAndDispatchIncidents(queue *queue.Queue, scrappedStatusPageIncidents []StatusPageIncident) error {
	incidentUpdates, err := storeStatusPageIncidents(scrappedStatusPageIncidents)
	if err != nil {
		return err
	}

	if queue != nil {
		publishUpdatesToDispatcher(queue, incidentUpdates)
	}

	return nil
}
func storeStatusPageIncidents(scrappedStatusPageIncidents []StatusPageIncident) ([]schema.IncidentUpdate, error) {
	statusPagesIncident, err := fetchAndStoreIncidents(scrappedStatusPageIncidents)
	if err != nil {
		return nil, err
	}

	incidentUpdates, incidentComponents := aggregateIncidentUpdatesAndIncidentComponents(statusPagesIncident)
	newIncidentUpdates := make([]schema.IncidentUpdate, 0)

	if len(incidentUpdates) > 0 {
		newIncidentUpdates, err = handleIncidentUpdates(incidentUpdates)
		if err != nil {
			return nil, err
		}
	}

	if len(incidentComponents) > 0 {
		err := handleIncidentComponents(incidentComponents, scrappedStatusPageIncidents[0].Incident.ServiceID)
		if err != nil {
			return nil, err
		}
	}

	return newIncidentUpdates, nil
}

// fetchAndStoreIncidents fetches existing incidents from DB and store
// new incident to the DB, while adding incident ID to []statusPageIncidents
func fetchAndStoreIncidents(statusPageIncidents []StatusPageIncident) ([]StatusPageIncident, error) {
	providerIDs := make([]string, 0)
	statusPageIncidentProviderIDMap := make(map[string]StatusPageIncident, 0)
	updatedStatusPageIncident := make([]StatusPageIncident, 0)

	for _, incident := range statusPageIncidents {
		providerIDs = append(providerIDs, incident.Incident.ProviderID)
		statusPageIncidentProviderIDMap[incident.Incident.ProviderID] = incident
	}

	existingIncidents, err := domain.Incident.GetByProviderIDs(providerIDs)
	if err != nil {
		return nil, err
	}

	existingIncidentsMap := lo.Associate(existingIncidents, func(incident schema.Incident) (string, schema.Incident) {
		return incident.ProviderID, incident
	})

	newIncidents := make([]schema.Incident, 0)

	for _, incident := range statusPageIncidents {
		existingIncident, ok := existingIncidentsMap[incident.Incident.ProviderID]
		if ok {
			statusPageIncidentFromProvider := statusPageIncidentProviderIDMap[incident.Incident.ProviderID]
			updatedStatusPageIncident = append(updatedStatusPageIncident, addIncidentIDToStatusPageIncident(statusPageIncidentFromProvider, existingIncident.BaseModel))
		} else {
			newIncidents = append(newIncidents, incident.Incident)
		}
	}

	insertedIncidents, err := domain.Incident.Create(newIncidents)
	if err != nil {
		return nil, err
	}

	for _, incident := range insertedIncidents {
		statusPageIncidentFromProvider, ok := statusPageIncidentProviderIDMap[incident.ProviderID]
		if ok {
			updatedStatusPageIncident = append(updatedStatusPageIncident, addIncidentIDToStatusPageIncident(statusPageIncidentFromProvider, incident.BaseModel))
		}
	}

	return updatedStatusPageIncident, nil

}

// addIncidentIDToStatusPageIncident populates incident ID to status page incident
// which inturn populates incident updates and incident components
// since it is currently empty
func addIncidentIDToStatusPageIncident(statusPageIncident StatusPageIncident, incidentBaseModel schema.BaseModel) StatusPageIncident {
	incidentUpdates := lo.Map(statusPageIncident.IncidentUpdates, func(incidentUpdate schema.IncidentUpdate, _ int) schema.IncidentUpdate {
		return schema.IncidentUpdate{
			IncidentID:     incidentBaseModel.ID,
			Description:    incidentUpdate.Description,
			ProviderID:     incidentUpdate.ProviderID,
			Status:         incidentUpdate.Status,
			ProviderStatus: incidentUpdate.ProviderStatus,
			StatusTime:     incidentUpdate.StatusTime,
		}
	})

	incidentComponents := lo.Map(statusPageIncident.IncidentComponents, func(incidentComponent StatusPageIncidentComponent, _ int) StatusPageIncidentComponent {
		return StatusPageIncidentComponent{
			IncidentComponents: schema.IncidentComponent{
				IncidentID:  incidentBaseModel.ID,
				ComponentID: incidentComponent.IncidentComponents.ComponentID,
			},
			ProviderComponentID: incidentComponent.ProviderComponentID,
			ComponentName:       incidentComponent.ComponentName,
		}

	})

	incident := schema.Incident{
		BaseModel:         incidentBaseModel,
		Name:              statusPageIncident.Incident.Name,
		Link:              statusPageIncident.Incident.Link,
		ProviderImpact:    statusPageIncident.Incident.ProviderImpact,
		Impact:            statusPageIncident.Incident.Impact,
		ServiceID:         statusPageIncident.Incident.ServiceID,
		ProviderID:        statusPageIncident.Incident.ProviderID,
		ProviderCreatedAt: statusPageIncident.Incident.ProviderCreatedAt,
	}

	return StatusPageIncident{
		Incident:           incident,
		IncidentUpdates:    incidentUpdates,
		IncidentComponents: incidentComponents,
	}
}

// aggregateIncidentUpdatesAndIncidentComponents denormalises Incident Updates
// and Incident Components from StatusPages incidents
// since they can be independently be inserted to the DB
func aggregateIncidentUpdatesAndIncidentComponents(statusPageIncidents []StatusPageIncident) ([]schema.IncidentUpdate, []StatusPageIncidentComponent) {
	incidentUpdates := make([]schema.IncidentUpdate, 0)
	incidentComponents := make([]StatusPageIncidentComponent, 0)
	for _, statusPageIncident := range statusPageIncidents {
		incidentUpdates = append(incidentUpdates, statusPageIncident.IncidentUpdates...)
		incidentComponents = append(incidentComponents, statusPageIncident.IncidentComponents...)
	}
	return incidentUpdates, incidentComponents
}

func handleIncidentUpdates(incidentUpdates []schema.IncidentUpdate) ([]schema.IncidentUpdate, error) {
	providerIDs := make([]string, 0)
	newToBeInsertedUpdates := make([]schema.IncidentUpdate, 0)

	for _, incidentUpdate := range incidentUpdates {
		providerIDs = append(providerIDs, incidentUpdate.ProviderID)
	}

	// fetch existing incident updates
	existingIncidentUpdates, err := domain.Incident.GetIncidentUpdatesByProviderIDs(providerIDs)
	if err != nil {
		return nil, err
	}

	existingIncidentUpdatesProviderIDMap := lo.Associate(existingIncidentUpdates, func(incidentUpdate schema.IncidentUpdate) (string, schema.IncidentUpdate) {
		return incidentUpdate.ProviderID, incidentUpdate
	})

	// insert new incident updates
	for _, incidentUpdate := range incidentUpdates {
		_, ok := existingIncidentUpdatesProviderIDMap[incidentUpdate.ProviderID]
		//Discard any existing incident update
		if !ok {
			newToBeInsertedUpdates = append(newToBeInsertedUpdates, incidentUpdate)
		}
	}

	insertedIncidentUpdates, err := domain.Incident.CreateIncidentUpdates(newToBeInsertedUpdates)
	if err != nil {
		return nil, err
	}

	return insertedIncidentUpdates, nil

}

func handleIncidentComponents(statusPageIncidentComponents []StatusPageIncidentComponent, serviceID uint) error {
	componentsByServiceID, err := domain.Component.GetAllByServiceID(serviceID)
	if err != nil {
		return err
	}
	componentsMap := lo.Associate(componentsByServiceID, func(component schema.Component) (string, schema.Component) {
		return component.ProviderID, component
	})

	var incidentComponents []schema.IncidentComponent
	for _, statusPageIncidentComponent := range statusPageIncidentComponents {
		component, ok := componentsMap[statusPageIncidentComponent.ProviderComponentID]
		if !ok {
			returnedComponents, err := domain.Component.Create([]schema.Component{{
				Name:       statusPageIncidentComponent.ComponentName,
				ProviderID: statusPageIncidentComponent.ProviderComponentID,
				ServiceID:  serviceID,
			}})

			if err != nil {
				return err
			}
			// Only one component will be returned
			// as we're inserting only one component
			component = returnedComponents[0]
		}

		incidentComponents = append(incidentComponents, schema.IncidentComponent{
			IncidentID:  statusPageIncidentComponent.IncidentComponents.IncidentID,
			ComponentID: component.ID,
		})
	}
	return domain.Incident.CreateIncidentComponents(incidentComponents)
}

func publishUpdatesToDispatcher(dispatcherQueue *queue.Queue, incidentUpdates []schema.IncidentUpdate) {
	for _, incidentUpdate := range incidentUpdates {
		dispatcherQueue.Publish(queue.IncidentPayload{
			State:          fmt.Sprintf("%v.%v", types.Incident, incidentUpdate.Status),
			IncidentUpdate: incidentUpdate,
		})
	}
}
