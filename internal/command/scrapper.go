package command

import (
	"context"
	"log/slog"
	"time"

	"github.com/yash492/statusy/internal/common/nullable"
	"github.com/yash492/statusy/internal/domain/components"
	"github.com/yash492/statusy/internal/domain/incidents"
	"github.com/yash492/statusy/internal/domain/services"
	"github.com/yash492/statusy/internal/domain/statuspage"
)

type ScrapperCmd struct {
	RegisteredStatuspages  map[string]statuspage.StatusPageProvider
	ServicesRepo           services.Repository
	IncidentsRepo          incidents.Repository
	IncidentUpdatesRepo    incidents.UpdatesRepository
	IncidentComponentsRepo incidents.ComponentsRepository
	ComponentsRepo         components.Repository
	ComponentGroupsRepo    components.GroupRepository
	logger                 *slog.Logger
}

type ScrapperParams struct {
	Services []services.ServiceResult
}

func (s ScrapperCmd) Execute(ctx context.Context, params ScrapperParams) error {
	start := time.Now()
	s.logger.InfoContext(ctx, "scrape cycle started", slog.Int("services_count", len(params.Services)))

	registeredServices := s.buildProviders(params.Services)
	s.logger.InfoContext(ctx, "providers prepared", slog.Int("providers_count", len(registeredServices)))
	scrappedComponents := []components.AggregateComponents{}
	scrappedIncidents := []incidents.Incident{}

	for _, service := range registeredServices {
		collectorStart := time.Now()
		s.logger.InfoContext(ctx, "scraping service", slog.String("slug", service.Slug().String()), slog.Uint64("service_id", uint64(service.ID())))

		componentsStart := time.Now()
		scrappedComponentForService, err := service.ScrapComponents()
		s.logger.InfoContext(
			ctx,
			"service components scraped",
			slog.String("slug", service.Slug().String()),
			slog.Duration("duration", time.Since(componentsStart)),
		)
		if err != nil {
			s.logger.ErrorContext(ctx, "error while scrapping components", slog.String("slug", service.Slug().String()), slog.Any("error", err))
		}

		incidentsStart := time.Now()
		scrappedIncidentForService, err := service.ScrapIncidents()
		s.logger.InfoContext(
			ctx,
			"service incidents scraped",
			slog.String("slug", service.Slug().String()),
			slog.Duration("duration", time.Since(incidentsStart)),
		)
		if err != nil {
			s.logger.ErrorContext(ctx, "error while scrapping incidents", slog.String("slug", service.Slug().String()), slog.Any("error", err))
		}

		scrappedComponentForService.Service.ID = service.ID()
		scrappedComponents = append(scrappedComponents, scrappedComponentForService)
		for _, incident := range scrappedIncidentForService {
			incident.ServiceID = service.ID()
			scrappedIncidents = append(scrappedIncidents, incident)
		}

		s.logger.InfoContext(
			ctx,
			"service scrape completed",
			slog.String("slug", service.Slug().String()),
			slog.Duration("duration", time.Since(collectorStart)),
		)
	}

	s.logger.InfoContext(
		ctx,
		"scrape collection completed",
		slog.Int("components_aggregate_count", len(scrappedComponents)),
		slog.Int("incidents_count", len(scrappedIncidents)),
	)

	componentsProviderMap, err := s.saveComponents(ctx, scrappedComponents)
	if err != nil {
		return err
	}

	err = s.saveIncidents(ctx, scrappedIncidents, componentsProviderMap)
	if err != nil {
		return err
	}

	s.logger.InfoContext(
		ctx,
		"scrape cycle completed",
		slog.Int("component_provider_map_count", len(componentsProviderMap)),
		slog.Duration("duration", time.Since(start)),
	)

	return nil
}

func (s ScrapperCmd) saveIncidents(
	ctx context.Context,
	scrappedIncidents []incidents.Incident,
	componentsProviderMap map[string]uint) error {
	s.logger.InfoContext(
		ctx,
		"saving incidents",
		slog.Int("incidents_count", len(scrappedIncidents)),
		slog.Int("components_provider_map_count", len(componentsProviderMap)),
	)

	incidentParams := make([]incidents.IncidentParams, 0, len(scrappedIncidents))
	for _, incident := range scrappedIncidents {
		incidentParams = append(incidentParams, incidents.IncidentParams{
			Title:             incident.Name,
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
	s.logger.InfoContext(ctx, "incidents saved", slog.Int("saved_incidents_count", len(savedIncidents)))

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
	s.logger.InfoContext(ctx, "incident updates saved", slog.Int("updates_count", len(updateParams)))

	_, err = s.IncidentComponentsRepo.SaveAll(ctx, componentParams)
	if err == nil {
		s.logger.InfoContext(ctx, "incident components saved", slog.Int("incident_components_count", len(componentParams)))
	}
	return err
}

func (s ScrapperCmd) saveComponents(ctx context.Context, scrappedComponents []components.AggregateComponents) (map[string]uint, error) {
	s.logger.InfoContext(ctx, "saving components", slog.Int("aggregate_components_count", len(scrappedComponents)))

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
	s.logger.InfoContext(ctx, "component groups saved", slog.Int("component_groups_count", len(componentGroupsResult)))
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
					ComponentGroupID: nullable.SetValue(providerComponentGroup.ID, true),
				})

			}
		}

		for _, component := range scrappedComponent.UngroupedComponents {
			componentGroup, ok := serviceComponentGroupMap[scrappedComponent.Service.ID][component.ProviderID]
			var componentGroupID nullable.Nullable[uint]
			if ok {
				componentGroupID = nullable.SetValue(componentGroup.ID, true)
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
	s.logger.InfoContext(ctx, "components saved", slog.Int("components_count", len(savedComponents)))

	// Component Provider ID -> Inserted Component ID
	componentProviderMap := map[string]uint{}
	for _, component := range savedComponents {
		componentProviderMap[component.ProviderID] = component.ID
	}

	return componentProviderMap, nil

}

func (s ScrapperCmd) buildProviders(servicesResult []services.ServiceResult) []statuspage.StatusPageProvider {
	servicesToBeScraped := make([]statuspage.StatusPageProvider, 0, len(servicesResult))

	for _, service := range servicesResult {
		provider, ok := s.RegisteredStatuspages[service.Slug]
		if !ok {
			s.logger.Warn("unsupported service slug", slog.String("slug", service.Slug))
			continue
		}

		servicesToBeScraped = append(servicesToBeScraped, provider.NewWithServiceID(service.ID))
	}

	s.logger.Info("providers selected for scrape", slog.Int("selected_count", len(servicesToBeScraped)))

	return servicesToBeScraped
}
