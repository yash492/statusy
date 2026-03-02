package atlassian

import (
	"github.com/samber/lo"
	"github.com/yash492/statusy/internal/common/nullable"
	"github.com/yash492/statusy/internal/domain/components"
	"github.com/yash492/statusy/internal/domain/incidents"
)

func FetchComponentsHelper(req ComponentsReq) components.AggregateComponents {

	ungroupedComponents := []components.Component{}
	groupedComponentsIDNameMap := map[string]string{}
	componentsToBeGrouped := []Component{}
	groupedComponents := map[string]components.ComponentGroup{}

	for _, atlassianComponent := range req.Components {
		if atlassianComponent.GroupID == nil {
			if !atlassianComponent.Group {
				ungroupedComponents = append(ungroupedComponents, components.Component{
					Name:       atlassianComponent.Name,
					ProviderID: atlassianComponent.ID,
				})

			} else if atlassianComponent.Group {
				groupedComponentsIDNameMap[atlassianComponent.ID] = atlassianComponent.Name
			}
		} else {
			componentsToBeGrouped = append(componentsToBeGrouped, atlassianComponent)
		}
	}

	for _, altassianComponent := range componentsToBeGrouped {
		if altassianComponent.GroupID == nil {
			continue
		}
		componentGroupName := groupedComponentsIDNameMap[*altassianComponent.GroupID]
		componentGroupID := *altassianComponent.GroupID
		groupedComponent, ok := groupedComponents[*altassianComponent.GroupID]
		if !ok {
			groupedComponent = components.ComponentGroup{
				Name:       componentGroupName,
				ProviderID: componentGroupID,
				Components: []components.Component{},
			}
		}
		component := components.Component{
			Name:       altassianComponent.Name,
			ProviderID: altassianComponent.ID,
		}

		groupedComponent.Components = append(groupedComponent.Components, component)
		groupedComponents[*altassianComponent.GroupID] = groupedComponent

	}

	return components.AggregateComponents{
		UngroupedComponents: ungroupedComponents,
		GroupedComponents:   lo.Values(groupedComponents),
	}
}

func FetchIncidentsHelper(req IncidentReq) []incidents.Incident {
	incidentList := []incidents.Incident{}
	for _, incidentReq := range req.Incidents {
		incident := incidents.Incident{
			Name:              incidentReq.Name,
			Link:              incidentReq.Shortlink,
			ProviderImpact:    nullable.SetValue(incidentReq.Impact),
			Impact:            nullable.SetValue(incidentReq.Impact),
			ProviderID:        incidentReq.ID,
			ProviderCreatedAt: incidentReq.CreatedAt,
		}

		incident.Updates = lo.Map(incidentReq.IncidentUpdates, func(update IncidentUpdate, _ int) incidents.IncidentUpdate {
			return incidents.IncidentUpdate{
				Description:        update.Body,
				IncidentProviderID: incident.ProviderID,
				ProviderID:         update.ID,
				Status:             update.Status,
				ProviderStatus:     update.Status,
				StatusTime:         update.CreatedAt,
			}
		})

		incident.Components = lo.Map(incidentReq.IncidentComponents, func(component IncidentComponent, _ int) components.Component {
			return components.Component{
				Name:       component.Name,
				ProviderID: component.ID,
			}
		})
		incidentList = append(incidentList, incident)
	}

	return incidentList
}
