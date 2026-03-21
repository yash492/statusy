package scheduledmaintenancesdb

import (
	"context"
	_ "embed"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/internal/domain/scheduledmaintenance"
)

//go:embed queries/get_scheduled_maintenances_by_service.sql
var getScheduledMaintenancesByServiceQuery string

type scheduledMaintenanceByServiceDto struct {
	ID                uint      `db:"id"`
	ServiceID         uint      `db:"service_id"`
	Title             string    `db:"title"`
	Status            string    `db:"status"`
	StartsAt          time.Time `db:"starts_at"`
	EndsAt            time.Time `db:"ends_at"`
	Link              string    `db:"link"`
	ProviderCreatedAt time.Time `db:"provider_created_at"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
}

func (c *PostgresScheduledMaintenanceRepository) GetByService(ctx context.Context, params scheduledmaintenance.ScheduledMaintenanceByServiceParams) ([]scheduledmaintenance.ScheduledMaintenanceByServiceResult, error) {
	rows, err := c.readDB.Query(ctx, getScheduledMaintenancesByServiceQuery, pgx.NamedArgs{
		"service_id": params.ServiceID,
		"limit":      params.Limit,
		"offset":     params.Offset,
	})
	if err != nil {
		c.lg.ErrorContext(ctx, "error querying scheduled maintenances by service", slog.Any("service_id", params.ServiceID), slog.Any("err", err))

		return nil, err
	}
	defer rows.Close()

	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[scheduledMaintenanceByServiceDto])
	if err != nil {
		if c.lg != nil {
			c.lg.ErrorContext(ctx, "error collecting scheduled maintenances by service rows", slog.Any("service_id", params.ServiceID), slog.Any("err", err))
		}
		return nil, err
	}

	results := make([]scheduledmaintenance.ScheduledMaintenanceByServiceResult, 0, len(dtos))
	for _, item := range dtos {
		results = append(results, scheduledmaintenance.ScheduledMaintenanceByServiceResult{
			ID:                item.ID,
			ServiceID:         item.ServiceID,
			Title:             item.Title,
			Status:            item.Status,
			StartsAt:          item.StartsAt,
			EndsAt:            item.EndsAt,
			ProviderCreatedAt: item.ProviderCreatedAt,
			Link:              item.Link,
		})
	}

	return results, nil
}
