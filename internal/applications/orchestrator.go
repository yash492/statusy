package applications

import (
	"context"
	"log/slog"

	"github.com/goccy/go-yaml"
	"github.com/yash492/statusy/internal/adapter/atlassianstatuspage"
	"github.com/yash492/statusy/internal/common"
	"github.com/yash492/statusy/internal/domain/components"
	"github.com/yash492/statusy/internal/domain/incidents"
	"github.com/yash492/statusy/internal/domain/services"
	"github.com/yash492/statusy/internal/domain/statuspage"
	"resty.dev/v3"
)

type ScrapperOrchestrator struct {
	ServicesYaml        []byte
	ServicesRepo        services.Repository
	IncidentsRepo       incidents.Repository
	ComponentsRepo      components.Repository
	ComponentGroupsRepo components.GroupRepository
	logger              *slog.Logger
}

type providerBuilder func(service services.ServiceResult) statuspage.StatusPageProvider

var providerBuilders = map[string]providerBuilder{
	"circleci": func(service services.ServiceResult) statuspage.StatusPageProvider {
		return atlassianstatuspage.NewCircleCIProvider(
			service.ComponentsUrl,
			service.IncidentsUrl,
			service.ScheduleMaintenancesUrl,
			service.ID,
			resty.New(),
		)
	},
}

func (s *ScrapperOrchestrator) Orchestrate() error {
	ctx := context.Background()

	var serviceParams []services.ServiceParams
	err := yaml.UnmarshalContext(ctx, s.ServicesYaml, &serviceParams)
	if err != nil {
		s.logger.ErrorContext(ctx, "error unmarshalling services yaml file", slog.Any("err", err))
		return err
	}

	servicesResult, err := s.ServicesRepo.SaveAll(ctx, serviceParams)
	if err != nil {
		return err
	}

	registeredServices := s.buildProviders(servicesResult)
	scrappedComponents := []components.AggregateComponents{}
	scrappedIncidents := []incidents.Incident{}

	for _, service := range registeredServices {
		scrappedComponentForService, err := service.ScrapComponents()
		if err != nil {
			s.logger.ErrorContext(ctx, "error while scrapping components for service %s", service.Slug(), slog.Any("error", err))
		}

		scrappedIncidentForService, err := service.ScrapIncidents()
		if err != nil {
			s.logger.ErrorContext(ctx, "error while scrapping components for service %s", service.Slug(), slog.Any("error", err))
		}

		scrappedComponents = append(scrappedComponents, scrappedComponentForService)
		scrappedIncidents = append(scrappedIncidents, scrappedIncidentForService...)
	}

	componentsProviderMap, err := s.saveComponents(ctx, scrappedComponents)
	if err != nil {
		return err
	}

	err = s.saveIncidents(ctx, scrappedIncidents, componentsProviderMap)
	if err != nil {
		return err
	}

	return nil
}

func (s *ScrapperOrchestrator) saveIncidents(
	ctx context.Context,
	scrappedIncidents []incidents.Incident,
	componentsProviderMap map[string]uint) error {

	incidentParams := make([]incidents.IncidentParams, 0, len(scrappedIncidents))
	for _, incident := range scrappedIncidents {
		incidentParams = append(incidentParams, incidents.IncidentParams{
			Name:              incident.Name,
			Link:              incident.Link,
			ProviderImpact:    incident.ProviderImpact,
			Impact:            incident.Impact,
			ServiceID:         incident.ServiceID,
			ProviderID:        incident.ProviderID,
			ProviderCreatedAt: incident.ProviderCreatedAt,
		})
	}

	savedIncidents, err := s.IncidentsRepo.SaveAll(ctx, incidentParams)
	if err != nil {
		return err
	}

	// TODO: save incident_updates using savedIncidents[i].ID + incident.Updates
	// TODO: save incident_components using savedIncidents[i].ID + componentsProviderMap to resolve component IDs
	_ = savedIncidents
	_ = componentsProviderMap

	return nil
}

func (s *ScrapperOrchestrator) saveComponents(ctx context.Context, scrappedComponents []components.AggregateComponents) (map[string]uint, error) {

	componentGroupsToBeScraped := make([]components.GroupParams, 0)
	for _, component := range scrappedComponents {
		for _, componentGroup := range component.GroupedComponents {
			componentGroupsToBeScraped = append(componentGroupsToBeScraped, components.GroupParams{
				Name:       componentGroup.Name,
				ProviderID: componentGroup.ProviderID,
				ServiceID:  component.Service.ID,
			})
		}
	}

	componentGroupsResult, err := s.ComponentGroupsRepo.SaveAll(ctx, componentGroupsToBeScraped)
	if err != nil {
		return nil, err
	}
	// Service ID -> ComponentGroup Provider ID -> ComponentGroup Result
	serviceComponentGroupMap := map[uint]map[string]components.ComponentGroupResult{}

	for _, componentGroup := range componentGroupsResult {
		componentgroupMap, ok := serviceComponentGroupMap[componentGroup.ServiceID]
		if !ok {
			componentgroupMap = make(map[string]components.ComponentGroupResult)
		}
		componentgroupMap[componentGroup.ProviderID] = componentGroup
		serviceComponentGroupMap[componentGroup.ServiceID] = componentgroupMap
	}

	componentParams := []components.ComponentParams{}

	for _, scrappedComponent := range scrappedComponents {
		for _, groupcomponent := range scrappedComponent.GroupedComponents {
			for _, component := range groupcomponent.Components {
				componentGroup := serviceComponentGroupMap[scrappedComponent.Service.ID][groupcomponent.ProviderID]

				componentParams = append(componentParams, components.ComponentParams{
					Name:       component.Name,
					ProviderID: component.ProviderID,
					ServiceID:  scrappedComponent.Service.ID,
					ComponentGroupID: common.Nullable[uint]{
						Value: componentGroup.ID,
						Valid: true,
					},
				})

			}
		}

		for _, component := range scrappedComponent.UngroupedComponents {
			componentGroup, ok := serviceComponentGroupMap[scrappedComponent.Service.ID][component.ProviderID]
			var componentGroupID common.Nullable[uint]
			if ok {
				componentGroupID = common.Nullable[uint]{
					Value: componentGroup.ID,
					Valid: true,
				}
			}

			componentParams = append(componentParams, components.ComponentParams{
				Name:             component.Name,
				ProviderID:       component.ProviderID,
				ServiceID:        scrappedComponent.Service.ID,
				ComponentGroupID: componentGroupID,
			})
		}
	}

	savedComponents, err := s.ComponentsRepo.SaveAll(ctx, componentParams)
	if err != nil {
		return nil, err
	}

	// Component Provider ID -> Inserted Component ID
	componentProviderMap := map[string]uint{}
	for _, component := range savedComponents {
		componentProviderMap[component.ProviderID] = component.ID
	}

	return componentProviderMap, nil

}

func (s *ScrapperOrchestrator) buildProviders(servicesResult []services.ServiceResult) []statuspage.StatusPageProvider {
	servicesToBeScraped := make([]statuspage.StatusPageProvider, 0, len(servicesResult))

	for _, service := range servicesResult {
		builder, ok := providerBuilders[service.Slug]
		if !ok {
			s.logger.Warn("unsupported service slug", slog.String("slug", service.Slug))
			continue
		}

		servicesToBeScraped = append(servicesToBeScraped, builder(service))
	}

	return servicesToBeScraped
}
