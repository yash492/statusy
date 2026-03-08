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
)

func (t *TestSuite) TestOrchestrate() {
	registeredStatusPage := collector.RegisterAll()

	orchestrator := &ScrapperCmd{
		RegisteredStatuspages:  registeredStatusPage,
		ServicesRepo:           servicesdb.NewPostgresServiceRepository(t.Logger, t.TestDb, t.TestDb),
		IncidentsRepo:          incidentsdb.NewPostgresIncidentRepository(t.Logger, t.TestDb, t.TestDb),
		IncidentUpdatesRepo:    incidentupdatesdb.NewPostgresIncidentUpdatesRepository(t.Logger, t.TestDb, t.TestDb),
		IncidentComponentsRepo: incidentcomponentsdb.NewPostgresIncidentComponentsRepository(t.Logger, t.TestDb, t.TestDb),
		ComponentsRepo:         componentsdb.NewPostgresComponentRepository(t.Logger, t.TestDb, t.TestDb),
		ComponentGroupsRepo:    componentgroupsdb.NewPostgresComponentGroupsRepository(t.Logger, t.TestDb, t.TestDb),
		logger:                 t.Logger,
	}

	err := orchestrator.Execute(context.Background())
	if err != nil {
		t.T().Fatalf("orchestrate failed: %s", err)
	}
}
