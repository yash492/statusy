package applications

import (
	"context"
	"log/slog"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yash492/statusy/internal/adapter/atlassianstatuspage"
	"github.com/yash492/statusy/internal/adapter/postgres/componentgroups"
	"github.com/yash492/statusy/internal/adapter/postgres/components"
	"github.com/yash492/statusy/internal/adapter/postgres/services"
	"github.com/yash492/statusy/internal/domain/statuspage"
	componentsGroupRepo "github.com/yash492/statusy/internal/repository/componentgroups"
	componentsRepo "github.com/yash492/statusy/internal/repository/components"

	servicesRepo "github.com/yash492/statusy/internal/repository/services"

	"resty.dev/v3"
)

type Deps struct {
	Logger  *slog.Logger
	ReadDB  *pgxpool.Pool
	WriteDB *pgxpool.Pool
}

type providerBuilder func(service servicesRepo.ServiceResult) statuspage.StatusPageProvider

var providerBuilders = map[string]providerBuilder{
	"circleci": func(service servicesRepo.ServiceResult) statuspage.StatusPageProvider {
		return atlassianstatuspage.NewCircleCIProvider(
			service.ComponentsUrl,
			service.IncidentsUrl,
			service.ScheduleMaintenancesUrl,
			service.ID,
			resty.New(),
		)
	},
}

func StartScrapper(deps Deps) {

	servicePostgresRepo := services.NewPostgresServiceRepository(
		deps.Logger,
		deps.ReadDB,
		deps.WriteDB,
	)

	componentPostgresRepo := components.NewPostgresComponentRepository(
		deps.Logger,
		deps.ReadDB,
		deps.WriteDB,
	)

	componentGroupPostgresRepo := componentgroups.NewPostgresComponentGroupsRepository(
		deps.Logger,
		deps.ReadDB,
		deps.WriteDB,
	)

	ctx := context.Background()

	servicesYaml, err := LoadServicesFromYaml(ctx, deps.Logger)
	if err != nil {
		deps.Logger.ErrorContext(ctx, "error reading services yaml file", slog.Any("err", err))
		return
	}

	var services []servicesRepo.ServiceParams
	err = yaml.UnmarshalContext(ctx, servicesYaml, &services)
	if err != nil {
		deps.Logger.ErrorContext(ctx, "error unmarshalling services yaml file", slog.Any("err", err))
		return
	}
	servicesResult, err := servicePostgresRepo.SaveAll(ctx, services)
	if err != nil {
		deps.Logger.ErrorContext(ctx, "error while saving services", slog.Any("err", err))
		return
	}

	servicesToBeScraped := buildProviders(servicesResult, deps.Logger)
	components := make([]statuspage.AggregateComponents, 0)

	//This part needs to be parallelized
	for _, service := range servicesToBeScraped {
		component, err := service.ScrapComponents()
		if err != nil {
			deps.Logger.ErrorContext(ctx, "error fetching components", slog.Any("err", err))
			continue
		}
		component.Service = statuspage.Service{
			ID:   service.ID(),
			Name: service.Name(),
		}
		components = append(components, component)
	}

	componentGroupsToBeScraped := make([]componentsGroupRepo.ComponentsGroupParams, 0)
	for _, component := range components {
		for _, componentGroup := range component.GroupedComponents {
			componentGroupsToBeScraped = append(componentGroupsToBeScraped, componentsGroupRepo.ComponentsGroupParams{
				Name:       componentGroup.Name,
				ProviderID: componentGroup.ProviderID,
				ServiceID:  component.Service.ID,
			})
		}
	}

	deps.Logger.InfoContext(ctx, "component_groups", slog.Any("component_groups", componentGroupsToBeScraped))

	componentGroupResult, err := componentGroupPostgresRepo.SaveAll(ctx, componentGroupsToBeScraped)
	if err != nil {
		deps.Logger.ErrorContext(ctx, "error saving component groups", slog.Any("err", err))
	}

	componentMap := map[uint]map[string]componentsGroupRepo.ComponentsGroupResult{}

	for _, componentGroup := range componentGroupResult {
		componentgroupMap, ok := componentMap[componentGroup.ServiceID]
		if !ok {
			componentgroupMap = make(map[string]componentsGroupRepo.ComponentsGroupResult)
		}
		componentgroupMap[componentGroup.ProviderID] = componentGroup
		componentMap[componentGroup.ServiceID] = componentgroupMap
	}

	componentsToBeSaved := []componentsRepo.ComponentParams{}

	for _, component := range components {
		for _, gcc := range component.GroupedComponents {
			for _, gc := range gcc.Components {
				componentGroup, ok := componentMap[component.Service.ID][gcc.ProviderID]
				var componentGroupID *uint
				if ok {
					componentGroupID = &componentGroup.ID
				}

				componentsToBeSaved = append(componentsToBeSaved, componentsRepo.ComponentParams{
					Name:             gc.Name,
					ProviderID:       gc.ProviderID,
					ServiceID:        component.Service.ID,
					ComponentGroupID: componentGroupID,
				})
			}
		}

		for _, gc := range component.UngroupedComponents {
			componentGroup, ok := componentMap[component.Service.ID][gc.ProviderID]
			var componentGroupID *uint
			if ok {
				componentGroupID = &componentGroup.ID
			}

			componentsToBeSaved = append(componentsToBeSaved, componentsRepo.ComponentParams{
				Name:             gc.Name,
				ProviderID:       gc.ProviderID,
				ServiceID:        component.Service.ID,
				ComponentGroupID: componentGroupID,
			})
		}
	}

	result, err := componentPostgresRepo.SaveAll(ctx, componentsToBeSaved)
	if err != nil {
		deps.Logger.ErrorContext(ctx, "error saving component ", slog.Any("err", err))
	}

	deps.Logger.InfoContext(ctx, "components", slog.Any("component", result))

}

func buildProviders(servicesResult []servicesRepo.ServiceResult, logger *slog.Logger) []statuspage.StatusPageProvider {
	servicesToBeScraped := make([]statuspage.StatusPageProvider, 0, len(servicesResult))

	for _, service := range servicesResult {
		builder, ok := providerBuilders[service.Slug]
		if !ok {
			logger.Warn("unsupported service slug", slog.String("slug", service.Slug))
			continue
		}

		servicesToBeScraped = append(servicesToBeScraped, builder(service))
	}

	return servicesToBeScraped
}

func LoadServicesFromYaml(ctx context.Context, lg *slog.Logger) ([]byte, error) {
	filePath := "../../data/services.yaml"
	yamlBytes, err := os.ReadFile(filePath)
	if err != nil {
		lg.ErrorContext(ctx, "error while reading sevices from yaml", slog.Any("err", err))
		return nil, err
	}

	return (yamlBytes), nil
}
