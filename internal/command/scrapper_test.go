package command

import (
	"context"

	"github.com/yash492/statusy/internal/adapter/collector"
	"github.com/yash492/statusy/internal/adapter/pgx/componentgroupsdb"
	"github.com/yash492/statusy/internal/adapter/pgx/componentsdb"
	"github.com/yash492/statusy/internal/adapter/pgx/incidentcomponentsdb"
	"github.com/yash492/statusy/internal/adapter/pgx/incidentsdb"
	"github.com/yash492/statusy/internal/adapter/pgx/incidentupdatesdb"
	"github.com/yash492/statusy/internal/adapter/pgx/servicesdb"
	"github.com/yash492/statusy/internal/domain/services"
)

func (t *TestSuite) TestOrchestrate() {

	ctx := context.Background()

	var serviceParams []services.ServiceParams

	registeredStatusPage := collector.RegisterAll()

	for _, provider := range registeredStatusPage {
		serviceParams = append(serviceParams, services.ServiceParams{
			Name: provider.Name(),
			Slug: provider.Slug().String(),
		})
	}

	servicesRepo := servicesdb.NewPostgresServiceRepository(t.Logger, t.TestDb, t.TestDb)

	servicesResult, err := servicesRepo.SaveAll(ctx, serviceParams)
	if err != nil {
		t.T().Fatalf("failed to save services: %s", err)
	}

	orchestrator := &ScrapperCmd{
		RegisteredStatuspages:  registeredStatusPage,
		ServicesRepo:           servicesRepo,
		IncidentsRepo:          incidentsdb.NewPostgresIncidentRepository(t.Logger, t.TestDb, t.TestDb),
		IncidentUpdatesRepo:    incidentupdatesdb.NewPostgresIncidentUpdatesRepository(t.Logger, t.TestDb, t.TestDb),
		IncidentComponentsRepo: incidentcomponentsdb.NewPostgresIncidentComponentsRepository(t.Logger, t.TestDb, t.TestDb),
		ComponentsRepo:         componentsdb.NewPostgresComponentRepository(t.Logger, t.TestDb, t.TestDb),
		ComponentGroupsRepo:    componentgroupsdb.NewPostgresComponentGroupsRepository(t.Logger, t.TestDb, t.TestDb),
		logger:                 t.Logger,
	}

	err = orchestrator.Execute(ctx, ScrapperParams{
		Services: servicesResult,
	})
	if err != nil {
		t.T().Fatalf("orchestrate failed: %s", err)
	}
}
