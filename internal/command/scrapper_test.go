package command

import (
	"context"

	"encoding/json"

	"github.com/yash492/statusy/internal/adapter/collector"
	"github.com/yash492/statusy/internal/adapter/pgx/componentgroupsdb"
	"github.com/yash492/statusy/internal/adapter/pgx/componentsdb"
	"github.com/yash492/statusy/internal/adapter/pgx/incidentcomponentsdb"
	"github.com/yash492/statusy/internal/adapter/pgx/incidentsdb"
	"github.com/yash492/statusy/internal/adapter/pgx/incidentupdatesdb"
	"github.com/yash492/statusy/internal/adapter/pgx/scheduledmaintenancesdb"
	"github.com/yash492/statusy/internal/adapter/pgx/scheduledmaintenanceupdatesdb"
	"github.com/yash492/statusy/internal/adapter/pgx/schedulemaintenancecomponentsdb"
	"github.com/yash492/statusy/internal/adapter/pgx/servicesdb"
	"github.com/yash492/statusy/internal/common/queue"
	"github.com/yash492/statusy/internal/domain/services"
	"resty.dev/v3"
)

type mockQueue struct{}

func (mockQueue) Send(ctx context.Context, queueName queue.QueueName, payload json.RawMessage) (string, error) {
	return "1", nil
}
func (mockQueue) SendBatch(ctx context.Context, queueName queue.QueueName, payloads []json.RawMessage) ([]string, error) {
	return []string{"1"}, nil
}
func (mockQueue) Read(ctx context.Context, queueName queue.QueueName, vt int, limit int) ([]queue.Message, error) {
	return nil, nil
}
func (mockQueue) Delete(ctx context.Context, queueName queue.QueueName, messageID string) error {
	return nil
}
func (mockQueue) Archive(ctx context.Context, queueName queue.QueueName, messageID string) error {
	return nil
}

func (t *TestSuite) TestOrchestrate() {

	ctx := context.Background()

	var serviceParams []services.ServiceParams

	registeredStatusPage := collector.RegisterAll(resty.New())

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

	incidentsRepo := incidentsdb.NewPostgresIncidentRepository(t.Logger, t.TestDb, t.TestDb)
	incidentUpdatesRepo := incidentupdatesdb.NewPostgresIncidentUpdatesRepository(t.Logger, t.TestDb, t.TestDb)
	incidentComponentsRepo := incidentcomponentsdb.NewPostgresIncidentComponentsRepository(t.Logger, t.TestDb, t.TestDb)
	scheduledMaintenanceRepo := scheduledmaintenancesdb.NewPostgresScheduledMaintenanceRepository(t.Logger, t.TestDb, t.TestDb)
	scheduledMaintenanceUpdatesRepo := scheduledmaintenanceupdatesdb.NewPostgresScheduledMaintenanceUpdatesRepository(t.Logger, t.TestDb, t.TestDb)
	scheduledMaintenanceComponentsRepo := schedulemaintenancecomponentsdb.NewPostgresScheduleMaintenanceComponentsRepository(t.Logger, t.TestDb, t.TestDb)
	componentsRepo := componentsdb.NewPostgresComponentRepository(t.Logger, t.TestDb, t.TestDb)
	componentGroupsRepo := componentgroupsdb.NewPostgresComponentGroupsRepository(t.Logger, t.TestDb, t.TestDb)

	logger := t.Logger

	orchestrator := NewScrapperCmd(
		registeredStatusPage,
		incidentsRepo,
		incidentUpdatesRepo,
		incidentComponentsRepo,
		scheduledMaintenanceRepo,
		scheduledMaintenanceUpdatesRepo,
		scheduledMaintenanceComponentsRepo,
		componentsRepo,
		componentGroupsRepo,
		mockQueue{},
		logger,
	)
	err = orchestrator.Execute(ctx, ScrapperParams{
		Services:    servicesResult,
		SkipEnqueue: true,
	})
	if err != nil {
		t.T().Fatalf("orchestrate failed: %s", err)
	}
}
