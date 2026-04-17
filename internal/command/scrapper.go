package command

import (
	"context"
	"log/slog"
	"time"

	"github.com/yash492/statusy/internal/common/nullable"
	"github.com/yash492/statusy/internal/domain/components"
	"github.com/yash492/statusy/internal/domain/incidents"
	"github.com/yash492/statusy/internal/domain/scheduledmaintenance"
	"github.com/yash492/statusy/internal/domain/services"
	"github.com/yash492/statusy/internal/domain/statuspage"
	"golang.org/x/sync/errgroup"
)

type ScrapperCmd struct {
	registeredStatuspages              map[string]statuspage.StatusPageProvider
	incidentsRepo                      incidents.Repository
	incidentUpdatesRepo                incidents.UpdatesRepository
	incidentComponentsRepo             incidents.ComponentsRepository
	scheduledMaintenanceRepo           scheduledmaintenance.Repository
	scheduledMaintenanceUpdatesRepo    scheduledmaintenance.UpdatesRepository
	scheduledMaintenanceComponentsRepo scheduledmaintenance.ComponentsRepository
	componentsRepo                     components.Repository
	componentGroupsRepo                components.GroupRepository
	logger                             *slog.Logger
}

func NewScrapperCmd(
	registeredStatuspages map[string]statuspage.StatusPageProvider,
	incidentsRepo incidents.Repository,
	incidentUpdatesRepo incidents.UpdatesRepository,
	incidentComponentsRepo incidents.ComponentsRepository,
	scheduledMaintenanceRepo scheduledmaintenance.Repository,
	scheduledMaintenanceUpdatesRepo scheduledmaintenance.UpdatesRepository,
	scheduledMaintenanceComponentsRepo scheduledmaintenance.ComponentsRepository,
	componentsRepo components.Repository,
	componentGroupsRepo components.GroupRepository,
	logger *slog.Logger,
) ScrapperCmd {
	return ScrapperCmd{
		registeredStatuspages:              registeredStatuspages,
		incidentsRepo:                      incidentsRepo,
		incidentUpdatesRepo:                incidentUpdatesRepo,
		incidentComponentsRepo:             incidentComponentsRepo,
		scheduledMaintenanceRepo:           scheduledMaintenanceRepo,
		scheduledMaintenanceUpdatesRepo:    scheduledMaintenanceUpdatesRepo,
		scheduledMaintenanceComponentsRepo: scheduledMaintenanceComponentsRepo,
		componentsRepo:                     componentsRepo,
		componentGroupsRepo:                componentGroupsRepo,
		logger:                             logger,
	}
}

type ScrapperParams struct {
	Services []services.ServiceResult
}

func (s ScrapperCmd) Execute(ctx context.Context, params ScrapperParams) error {
	start := time.Now()
	s.logger.InfoContext(ctx, "scrape cycle started", slog.Int("services_count", len(params.Services)))

	registeredServices := s.buildProviders(params.Services)
	s.logger.InfoContext(ctx, "providers prepared", slog.Int("providers_count", len(registeredServices)))
	type serviceScrapeResult struct {
		Component             components.AggregateComponents
		Incidents             []incidents.Incident
		ScheduledMaintenances []scheduledmaintenance.ScheduledMaintenance
	}

	scrappedComponents := make([]components.AggregateComponents, 0, len(registeredServices))
	scrappedIncidents := []incidents.Incident{}
	scrappedScheduledMaintenances := []scheduledmaintenance.ScheduledMaintenance{}

	scrapedResults := make([]serviceScrapeResult, len(registeredServices))
	limitGroup := new(errgroup.Group)
	limitGroup.SetLimit(10)

	for i, service := range registeredServices {
		i, service := i, service
		limitGroup.Go(func() error {
			collectorStart := time.Now()
			s.logger.InfoContext(ctx, "scraping service", slog.String("slug", service.Slug().String()), slog.Uint64("service_id", uint64(service.ID())))

			var subGroup errgroup.Group
			var sc components.AggregateComponents
			var si []incidents.Incident
			var ssm []scheduledmaintenance.ScheduledMaintenance

			subGroup.Go(func() error {
				componentsStart := time.Now()
				c, err := service.ScrapComponents()
				s.logger.InfoContext(
					ctx,
					"service components scraped",
					slog.String("slug", service.Slug().String()),
					slog.Int64("duration_ms", time.Since(componentsStart).Milliseconds()),
				)
				if err != nil {
					s.logger.ErrorContext(ctx, "error while scrapping components", slog.String("slug", service.Slug().String()), slog.Any("error", err))
				}
				sc = c
				return nil
			})

			subGroup.Go(func() error {
				incidentsStart := time.Now()
				inc, err := service.ScrapIncidents()
				s.logger.InfoContext(
					ctx,
					"service incidents scraped",
					slog.String("slug", service.Slug().String()),
					slog.Int64("duration_ms", time.Since(incidentsStart).Milliseconds()),
				)
				if err != nil {
					s.logger.ErrorContext(ctx, "error while scrapping incidents", slog.String("slug", service.Slug().String()), slog.Any("error", err))
				}
				si = inc
				return nil
			})

			subGroup.Go(func() error {
				scheduledMaintenanceStart := time.Now()
				sm, err := service.ScrapScheduledMaintenance()
				s.logger.InfoContext(
					ctx,
					"service scheduled maintenances scraped",
					slog.String("slug", service.Slug().String()),
					slog.Int64("duration_ms", time.Since(scheduledMaintenanceStart).Milliseconds()),
				)
				if err != nil {
					s.logger.ErrorContext(ctx, "error while scrapping scheduled maintenances", slog.String("slug", service.Slug().String()), slog.Any("error", err))
				}
				ssm = sm
				return nil
			})

			_ = subGroup.Wait()

			sc.Service.ID = service.ID()
			for idx := range si {
				si[idx].ServiceID = service.ID()
			}
			for idx := range ssm {
				ssm[idx].ServiceID = service.ID()
			}

			scrapedResults[i] = serviceScrapeResult{
				Component:             sc,
				Incidents:             si,
				ScheduledMaintenances: ssm,
			}

			s.logger.InfoContext(
				ctx,
				"service scrape completed",
				slog.String("slug", service.Slug().String()),
				slog.Int64("duration_ms", time.Since(collectorStart).Milliseconds()),
			)
			return nil
		})
	}

	_ = limitGroup.Wait()

	for _, result := range scrapedResults {
		scrappedComponents = append(scrappedComponents, result.Component)
		scrappedIncidents = append(scrappedIncidents, result.Incidents...)
		scrappedScheduledMaintenances = append(scrappedScheduledMaintenances, result.ScheduledMaintenances...)
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

	err = s.saveScheduledMaintenance(ctx, scrappedScheduledMaintenances, componentsProviderMap)
	if err != nil {
		return err
	}

	s.logger.InfoContext(
		ctx,
		"scrape cycle completed",
		slog.Int("component_provider_map_count", len(componentsProviderMap)),
		slog.Int64("duration_ms", time.Since(start).Milliseconds()),
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

	savedIncidents, err := s.incidentsRepo.SaveAll(ctx, incidentParams)
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

	_, err = s.incidentUpdatesRepo.SaveAll(ctx, updateParams)
	if err != nil {
		return err
	}
	s.logger.InfoContext(ctx, "incident updates saved", slog.Int("updates_count", len(updateParams)))

	_, err = s.incidentComponentsRepo.SaveAll(ctx, componentParams)
	if err == nil {
		s.logger.InfoContext(ctx, "incident components saved", slog.Int("incident_components_count", len(componentParams)))
	}
	return err
}

func (s ScrapperCmd) saveScheduledMaintenance(
	ctx context.Context,
	scrappedScheduledMaintenances []scheduledmaintenance.ScheduledMaintenance,
	componentsProviderMap map[string]uint) error {
	s.logger.InfoContext(
		ctx,
		"saving scheduled maintenances",
		slog.Int("scheduled_maintenances_count", len(scrappedScheduledMaintenances)),
		slog.Int("components_provider_map_count", len(componentsProviderMap)),
	)

	smParams := make([]scheduledmaintenance.ScheduledMaintenanceParams, 0, len(scrappedScheduledMaintenances))
	for _, sm := range scrappedScheduledMaintenances {
		smParams = append(smParams, scheduledmaintenance.ScheduledMaintenanceParams{
			Title:             sm.Name,
			Link:              sm.Link,
			ProviderImpact:    sm.ProviderImpact,
			Impact:            sm.Impact,
			StartsAt:          sm.StartsAt,
			EndsAt:            sm.EndsAt,
			ServiceID:         sm.ServiceID,
			ProviderID:        sm.ProviderID,
			ProviderCreatedAt: sm.ProviderCreatedAt,
		})
	}

	savedSMs, err := s.scheduledMaintenanceRepo.SaveAll(ctx, smParams)
	if err != nil {
		return err
	}
	s.logger.InfoContext(ctx, "scheduled maintenances saved", slog.Int("saved_count", len(savedSMs)))

	// Build a map from provider_id -> saved SM for quick lookup
	savedSMsByProviderID := make(map[string]scheduledmaintenance.ScheduledMaintenanceResult, len(savedSMs))
	for _, saved := range savedSMs {
		savedSMsByProviderID[saved.ProviderID] = saved
	}

	updateParams := make([]scheduledmaintenance.ScheduledMaintenanceUpdateParams, 0)
	componentParams := make([]scheduledmaintenance.ScheduledMaintenanceComponentParams, 0)

	for _, sm := range scrappedScheduledMaintenances {
		saved := savedSMsByProviderID[sm.ProviderID]

		for _, update := range sm.Updates {
			updateParams = append(updateParams, scheduledmaintenance.ScheduledMaintenanceUpdateParams{
				ScheduledMaintenanceID: saved.ID,
				Description:            update.Description,
				ProviderID:             update.ProviderID,
				ProviderStatus:         update.ProviderStatus,
				Status:                 update.Status,
				StatusTime:             update.StatusTime,
			})
		}

		for _, component := range sm.Components {
			componentID, ok := componentsProviderMap[component.ProviderID]
			if !ok {
				continue
			}
			componentParams = append(componentParams, scheduledmaintenance.ScheduledMaintenanceComponentParams{
				ScheduledMaintenanceID: saved.ID,
				ComponentID:            componentID,
			})
		}
	}

	_, err = s.scheduledMaintenanceUpdatesRepo.SaveAll(ctx, updateParams)
	if err != nil {
		return err
	}
	s.logger.InfoContext(ctx, "scheduled maintenance updates saved", slog.Int("updates_count", len(updateParams)))

	_, err = s.scheduledMaintenanceComponentsRepo.SaveAll(ctx, componentParams)
	if err == nil {
		s.logger.InfoContext(ctx, "scheduled maintenance components saved", slog.Int("components_count", len(componentParams)))
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

	componentGroupsResult, err := s.componentGroupsRepo.SaveAll(ctx, componentGroupsToBeScraped)
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

	savedComponents, err := s.componentsRepo.SaveAll(ctx, componentParams)
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
		provider, ok := s.registeredStatuspages[service.Slug]
		if !ok {
			s.logger.Warn("unsupported service slug", slog.String("slug", service.Slug))
			continue
		}

		servicesToBeScraped = append(servicesToBeScraped, provider.NewWithServiceID(service.ID))
	}

	s.logger.Info("providers selected for scrape", slog.Int("selected_count", len(servicesToBeScraped)))

	return servicesToBeScraped
}
