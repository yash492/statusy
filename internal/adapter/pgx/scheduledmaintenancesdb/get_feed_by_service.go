package scheduledmaintenancesdb

import (
	"context"
	_ "embed"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/internal/domain/scheduledmaintenance"
)

//go:embed queries/get_feed_sm_by_service.sql
var getFeedSMByServiceQuery string

type feedSMByServiceDto struct {
	ID                 uint      `db:"id"`
	ServiceID          uint      `db:"service_id"`
	Title              string    `db:"title"`
	Status             string    `db:"status"`
	Link               string    `db:"link"`
	ProviderCreatedAt  time.Time `db:"provider_created_at"`
	AffectedComponents string    `db:"affected_components"`
}

func (c *PostgresScheduledMaintenanceRepository) GetFeedByService(ctx context.Context, params scheduledmaintenance.ScheduledMaintenanceByServiceParams) ([]scheduledmaintenance.FeedScheduledMaintenanceByServiceResult, error) {
	rows, err := c.readDB.Query(ctx, getFeedSMByServiceQuery, pgx.NamedArgs{
		"service_id": params.ServiceID,
		"limit":      params.Limit,
		"offset":     params.Offset,
	})
	if err != nil {
		if c.lg != nil {
			c.lg.ErrorContext(ctx, "error querying feed SM by service", slog.Any("service_id", params.ServiceID), slog.Any("err", err))
		}
		return nil, err
	}
	defer rows.Close()

	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[feedSMByServiceDto])
	if err != nil {
		if c.lg != nil {
			c.lg.ErrorContext(ctx, "error collecting feed SM by service rows", slog.Any("service_id", params.ServiceID), slog.Any("err", err))
		}
		return nil, err
	}

	results := make([]scheduledmaintenance.FeedScheduledMaintenanceByServiceResult, 0, len(dtos))
	for _, item := range dtos {
		results = append(results, scheduledmaintenance.FeedScheduledMaintenanceByServiceResult{
			ID:                 item.ID,
			ServiceID:          item.ServiceID,
			Title:              item.Title,
			Status:             item.Status,
			ProviderCreatedAt:  item.ProviderCreatedAt,
			Link:               item.Link,
			AffectedComponents: item.AffectedComponents,
		})
	}

	return results, nil
}
