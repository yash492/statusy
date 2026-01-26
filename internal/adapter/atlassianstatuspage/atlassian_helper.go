package atlassianstatuspage

import (
	"github.com/samber/lo"
	"github.com/yash492/statusy/internal/common"
	"github.com/yash492/statusy/internal/domain/statuspage"
)

func fetchComponentsHelper(req atlassianComponentsReq) statuspage.AggregateComponents {

	ungroupedComponents := []statuspage.Component{}
	groupedComponentsIDNameMap := map[string]string{}
	componentsToBeGrouped := []atlassianComponent{}
	groupedComponents := map[string]statuspage.ComponentGroup{}

	for _, atlassianComponent := range req.Components {
		if atlassianComponent.GroupID == nil {
			if !atlassianComponent.Group {
				ungroupedComponents = append(ungroupedComponents, statuspage.Component{
					Name:        atlassianComponent.Name,
					ProviderID:  atlassianComponent.ID,
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
			groupedComponent = statuspage.ComponentGroup{
				Name:       componentGroupName,
				ProviderID: componentGroupID,
				Components: []statuspage.Component{},
			}
		}
		component := statuspage.Component{
			Name:        altassianComponent.Name,
			ProviderID:  altassianComponent.ID,
		}

		groupedComponent.Components = append(groupedComponent.Components, component)
		groupedComponents[*altassianComponent.GroupID] = groupedComponent

	}

	return statuspage.AggregateComponents{
		UngroupedComponents: ungroupedComponents,
		GroupedComponents:   lo.Values(groupedComponents),
	}
}

func fetchIncidentsHelper(req atlassianIncidentReq, serviceSlug string) []statuspage.Incident {
	incidents := []statuspage.Incident{}
	for _, incidentReq := range req.Incidents {
		incident := statuspage.Incident{
			Name:              incidentReq.Name,
			Link:              incidentReq.Shortlink,
			ServiceSlug:       serviceSlug,
			ProviderImpact:    common.SetNullableValue(incidentReq.Impact),
			Impact:            common.SetNullableValue(incidentReq.Impact),
			ProviderID:        incidentReq.ID,
			ProviderCreatedAt: incidentReq.CreatedAt,
		}

		incident.Updates = lo.Map(incidentReq.IncidentUpdates, func(update atlassianIncidentUpdate, _ int) statuspage.IncidentUpdate {
			return statuspage.IncidentUpdate{
				Description:        update.Body,
				IncidentProviderID: incident.ProviderID,
				ProviderID:         update.ID,
				Status:             update.Status,
				ProviderStatus:     update.Status,
				StatusTime:         update.CreatedAt,
			}
		})

		incident.Components = lo.Map(incidentReq.IncidentComponents, func(component atlassianIncidentComponent, _ int) statuspage.Component {
			return statuspage.Component{
				Name:        component.Name,
				ProviderID:  serviceSlug,
			}
		})
		incidents = append(incidents, incident)
	}

	return incidents
}