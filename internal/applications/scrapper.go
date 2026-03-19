package applications

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yash492/statusy/internal/adapter/collector"
	"github.com/yash492/statusy/internal/adapter/pgx/componentgroupsdb"
	"github.com/yash492/statusy/internal/adapter/pgx/componentsdb"
	"github.com/yash492/statusy/internal/adapter/pgx/incidentcomponentsdb"
	"github.com/yash492/statusy/internal/adapter/pgx/incidentsdb"
	"github.com/yash492/statusy/internal/adapter/pgx/incidentupdatesdb"
	"github.com/yash492/statusy/internal/adapter/pgx/servicesdb"
	"github.com/yash492/statusy/internal/command"
	"github.com/yash492/statusy/internal/domain/services"
	"github.com/yash492/statusy/internal/domain/statuspage"
	"resty.dev/v3"
)

type ScrapperDeps struct {
	lg      *slog.Logger
	readDB  *pgxpool.Pool
	writeDB *pgxpool.Pool
}

func NewScrapperDeps(
	lg *slog.Logger,
	readDB *pgxpool.Pool,
	writeDB *pgxpool.Pool,
) ScrapperDeps {
	return ScrapperDeps{
		lg:      lg,
		readDB:  readDB,
		writeDB: writeDB,
	}
}

type ScrapperApplication struct {
	cmd                  scrapperCommand
	servicesRepo         *servicesdb.PostgresServiceRepository
	registeredStatusPage map[string]statuspage.StatusPageProvider
	lg                   *slog.Logger
}

type scrapperCommand struct {
	scrapperCmd command.ScrapperCmd
}

func NewScrapperApplication(deps ScrapperDeps) ScrapperApplication {
	lg := deps.lg
	servicesRepo := servicesdb.NewPostgresServiceRepository(
		lg,
		deps.readDB,
		deps.writeDB,
	)
	incidentsRepo := incidentsdb.NewPostgresIncidentRepository(
		lg,
		deps.readDB,
		deps.writeDB,
	)
	incidentUpdatesRepo := incidentupdatesdb.NewPostgresIncidentUpdatesRepository(
		lg,
		deps.readDB,
		deps.writeDB,
	)
	incidentComponentsRepo := incidentcomponentsdb.NewPostgresIncidentComponentsRepository(
		lg,
		deps.readDB,
		deps.writeDB,
	)
	componentsRepo := componentsdb.NewPostgresComponentRepository(
		lg,
		deps.readDB,
		deps.writeDB,
	)
	componentGroupsRepo := componentgroupsdb.NewPostgresComponentGroupsRepository(
		lg,
		deps.readDB,
		deps.writeDB,
	)

	registeredStatusPage := collector.RegisterAll(resty.New())

	return ScrapperApplication{
		lg:                   deps.lg,
		registeredStatusPage: registeredStatusPage,
		servicesRepo:         servicesRepo,
		cmd: scrapperCommand{
			scrapperCmd: command.ScrapperCmd{
				RegisteredStatuspages:  registeredStatusPage,
				ServicesRepo:           servicesRepo,
				IncidentsRepo:          incidentsRepo,
				IncidentUpdatesRepo:    incidentUpdatesRepo,
				IncidentComponentsRepo: incidentComponentsRepo,
				ComponentsRepo:         componentsRepo,
				ComponentGroupsRepo:    componentGroupsRepo,
			},
		},
	}
}

func (s ScrapperApplication) Start(ctx context.Context, scrapInterval int) error {
	if scrapInterval <= 0 {
		return fmt.Errorf("invalid scrap interval: %d (must be > 0)", scrapInterval)
	}

	var serviceParams []services.ServiceParams

	for _, provider := range s.registeredStatusPage {
		serviceParams = append(serviceParams, services.ServiceParams{
			Name: provider.Name(),
			Slug: provider.Slug().String(),
		})
	}

	servicesResult, err := s.servicesRepo.SaveAll(ctx, serviceParams)
	if err != nil {
		return err
	}

	executeScrape := func() {
		err := s.cmd.scrapperCmd.Execute(ctx, command.ScrapperParams{
			Services: servicesResult,
		})

		if err != nil {
			s.lg.ErrorContext(ctx, "scraper command execution failed", "error", err)
		}
	}

	// Run once at startup so we don't wait an entire interval before first scrape.
	executeScrape()

	ticker := time.NewTicker(time.Duration(scrapInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.lg.InfoContext(ctx, "scraper stopped", "reason", ctx.Err())
			return ctx.Err()
		case <-ticker.C:
			executeScrape()
		}
	}

}
