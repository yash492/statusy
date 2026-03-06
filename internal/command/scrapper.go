package command

import (
	"context"
	"log/slog"

	"github.com/goccy/go-yaml"
	"github.com/yash492/statusy/internal/adapter/collector/registry"
	"github.com/yash492/statusy/internal/common/nullable"
	"github.com/yash492/statusy/internal/domain/components"
	"github.com/yash492/statusy/internal/domain/incidents"
	"github.com/yash492/statusy/internal/domain/services"
	"github.com/yash492/statusy/internal/domain/statuspage"
)

type ScrapperCmd struct {
	RegisteredStatuspages  map[string]registry.ProviderBuilderFunc
	ServicesYaml           []byte
	ServicesRepo           services.Repository
	IncidentsRepo          incidents.Repository
	IncidentUpdatesRepo    incidents.UpdatesRepository
	IncidentComponentsRepo incidents.ComponentsRepository
	ComponentsRepo         components.Repository
	ComponentGroupsRepo    components.GroupRepository
	logger                 *slog.Logger
}

func (s *ScrapperCmd) Execute(ctx context.Context) error {
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
			s.logger.ErrorContext(ctx, "error while scrapping components for service %s", string(service.Slug()), slog.Any("error", err))
		}

		scrappedIncidentForService, err := service.ScrapIncidents()
		if err != nil {
			s.logger.ErrorContext(ctx, "error while scrapping components for service %s", string(service.Slug()), slog.Any("error", err))
		}

		scrappedComponentForService.Service.ID = service.ID()
		scrappedComponents = append(scrappedComponents, scrappedComponentForService)
		for _, incident := range scrappedIncidentForService {
			incident.ServiceID = service.ID()
			scrappedIncidents = append(scrappedIncidents, incident)
		}
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

func (s *ScrapperCmd) saveIncidents(
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

	// Build a map from provider_id -> saved incident for quick lookup
	savedIncidentsByProviderID := make(map[string]incidents.IncidentResult, len(savedIncidents))
	for _, saved := range savedIncidents {
		savedIncidentsByProviderID[saved.ProviderID] = saved
	}

	updateParams := make([]incidents.IncidentUpdateParams, 0)
	componentParams := make([]incidents.IncidentComponentParams, 0)

	for _, incident := range scrappedIncidents {
		saved := savedIncidentsByProviderID[incident.ProviderID]

		for _, update := range incident.Updates {
			updateParams = append(updateParams, incidents.IncidentUpdateParams{
				IncidentID:     saved.ID,
				Description:    update.Description,
				ProviderID:     update.ProviderID,
				ProviderStatus: update.ProviderStatus,
				Status:         update.Status,
				StatusTime:     update.StatusTime,
			})
		}

		for _, component := range incident.Components {
			componentID, ok := componentsProviderMap[component.ProviderID]
			if !ok {
				continue
			}
			componentParams = append(componentParams, incidents.IncidentComponentParams{
				IncidentID:  saved.ID,
				ComponentID: componentID,
			})
		}
	}

	_, err = s.IncidentUpdatesRepo.SaveAll(ctx, updateParams)
	if err != nil {
		return err
	}

	_, err = s.IncidentComponentsRepo.SaveAll(ctx, componentParams)
	return err
}

func (s *ScrapperCmd) saveComponents(ctx context.Context, scrappedComponents []components.AggregateComponents) (map[string]uint, error) {

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
				providerComponentGroup := serviceComponentGroupMap[scrappedComponent.Service.ID][groupcomponent.ProviderID]

				componentParams = append(componentParams, components.ComponentParams{
					Name:             component.Name,
					ProviderID:       component.ProviderID,
					ServiceID:        scrappedComponent.Service.ID,
					ComponentGroupID: nullable.SetValue(providerComponentGroup.ID),
				})

			}
		}

		for _, component := range scrappedComponent.UngroupedComponents {
			componentGroup, ok := serviceComponentGroupMap[scrappedComponent.Service.ID][component.ProviderID]
			var componentGroupID nullable.Nullable[uint]
			if ok {
				componentGroupID = nullable.SetValue(componentGroup.ID)
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

func (s *ScrapperCmd) buildProviders(servicesResult []services.ServiceResult) []statuspage.StatusPageProvider {
	servicesToBeScraped := make([]statuspage.StatusPageProvider, 0, len(servicesResult))

	for _, service := range servicesResult {
		builder, ok := s.RegisteredStatuspages[service.Slug]
		if !ok {
			s.logger.Warn("unsupported service slug", slog.String("slug", service.Slug))
			continue
		}

		servicesToBeScraped = append(servicesToBeScraped, builder(service))
	}

	return servicesToBeScraped
}
