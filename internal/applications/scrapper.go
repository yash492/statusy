package applications

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	domaincomponents "github.com/yash492/statusy/internal/domain/components"
	domainservices "github.com/yash492/statusy/internal/domain/services"
)

type Deps struct {
	Logger  *slog.Logger
	ReadDB  *pgxpool.Pool
	WriteDB *pgxpool.Pool
}

func StartScrapper(deps Deps) {

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

		incidents, err := service.ScrapIncidents()
		if err != nil {
			deps.Logger.ErrorContext(ctx, "error fetching incidents", slog.Any("err", err))
			continue
		}
		for _, incident := range incidents {
			incident.ServiceID = service.ID()

			incidentsList = append(incidentsList)
		}
	}

	componentMap := map[uint]map[string]domaincomponents.ComponentGroupResult{}

	result, err := componentPostgresRepo.SaveAll(ctx, componentsToBeSaved)
	if err != nil {
		deps.Logger.ErrorContext(ctx, "error saving component ", slog.Any("err", err))
	}

	deps.Logger.InfoContext(ctx, "components", slog.Any("component", result))

	// incidentResult :=

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
