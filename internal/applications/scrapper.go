package applications

import (
	"context"
	"log/slog"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yash492/statusy/internal/adapter/atlassianstatuspage"
	"github.com/yash492/statusy/internal/adapter/pgx/componentgroupsdb"
	"github.com/yash492/statusy/internal/adapter/pgx/componentsdb"
	"github.com/yash492/statusy/internal/adapter/pgx/servicesdb"
	"github.com/yash492/statusy/internal/common"
	domaincomponents "github.com/yash492/statusy/internal/domain/components"
	domainservices "github.com/yash492/statusy/internal/domain/services"
	"github.com/yash492/statusy/internal/domain/statuspage"

	"resty.dev/v3"
)

type Deps struct {
	Logger  *slog.Logger
	ReadDB  *pgxpool.Pool
	WriteDB *pgxpool.Pool
}

type providerBuilder func(service domainservices.ServiceResult) statuspage.StatusPageProvider

var providerBuilders = map[string]providerBuilder{
	"circleci": func(service domainservices.ServiceResult) statuspage.StatusPageProvider {
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

	servicePostgresRepo := servicesdb.NewPostgresServiceRepository(
		deps.Logger,
		deps.ReadDB,
		deps.WriteDB,
	)

	componentPostgresRepo := componentsdb.NewPostgresComponentRepository(
		deps.Logger,
		deps.ReadDB,
		deps.WriteDB,
	)

	componentGroupPostgresRepo := componentgroupsdb.NewPostgresComponentGroupsRepository(
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

	var services []domainservices.ServiceParams
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
	componentsList := make([]domaincomponents.AggregateComponents, 0)

	//This part needs to be parallelized
	for _, service := range servicesToBeScraped {
		component, err := service.ScrapComponents()
		if err != nil {
			deps.Logger.ErrorContext(ctx, "error fetching components", slog.Any("err", err))
			continue
		}
		component.Service = domainservices.Service{
			ID:   service.ID(),
			Name: service.Name(),
		}
		componentsList = append(componentsList, component)
	}

	componentGroupsToBeScraped := make([]domaincomponents.GroupParams, 0)
	for _, component := range componentsList {
		for _, componentGroup := range component.GroupedComponents {
			componentGroupsToBeScraped = append(componentGroupsToBeScraped, domaincomponents.GroupParams{
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

	componentMap := map[uint]map[string]domaincomponents.ComponentResult{}

	for _, componentGroup := range componentGroupResult {
		componentgroupMap, ok := componentMap[componentGroup.ServiceID]
		if !ok {
			componentgroupMap = make(map[string]domaincomponents.Result)
		}
		componentgroupMap[componentGroup.ProviderID] = componentGroup
		componentMap[componentGroup.ServiceID] = componentgroupMap
	}

	componentsToBeSaved := []domaincomponents.ComponentParams{}

	for _, component := range componentsList {
		for _, gcc := range component.GroupedComponents {
			for _, gc := range gcc.Components {
				componentGroup := componentMap[component.Service.ID][gcc.ProviderID]

				componentsToBeSaved = append(componentsToBeSaved, domaincomponents.ComponentParams{
					Name:             gc.Name,
					ProviderID:       gc.ProviderID,
					ServiceID:        component.Service.ID,
					ComponentGroupID: common.Nullable[uint]{},
				})
			}
		}

		for _, gc := range component.UngroupedComponents {
			componentGroup, ok := componentMap[component.Service.ID][gc.ProviderID]
			var componentGroupID *uint
			if ok {
				componentGroupID = &componentGroup.ID
			}

			componentsToBeSaved = append(componentsToBeSaved, domaincomponents.ComponentParams{
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

func buildProviders(servicesResult []domainservices.ServiceResult, logger *slog.Logger) []statuspage.StatusPageProvider {
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
