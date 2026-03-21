package schedulemaintenancecomponentsdb

import (
	"context"
	_ "embed"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/samber/lo"
	"github.com/yash492/statusy/internal/common/nullable"
	"github.com/yash492/statusy/internal/domain/scheduledmaintenance"
)

//go:embed queries/insert_schedule_maintenance_components.sql
var insertScheduleMaintenanceComponentsQuery string

type scheduleMaintenanceComponentDto struct {
	ID                     uint             `db:"id"`
	ScheduledMaintenanceID uint             `db:"scheduled_maintenance_id"`
	ComponentID            uint             `db:"component_id"`
	CreatedAt              time.Time        `db:"created_at"`
	UpdatedAt              time.Time        `db:"updated_at"`
	DeletedAt              pgtype.Timestamp `db:"deleted_at"`
}

func (r *PostgresScheduleMaintenanceComponentsRepository) SaveAll(ctx context.Context, params []scheduledmaintenance.ScheduledMaintenanceComponentParams) ([]scheduledmaintenance.ScheduledMaintenanceComponentResult, error) {
	batchInserts := &pgx.Batch{}
	componentsResponse := []scheduleMaintenanceComponentDto{}

	for _, param := range params {
		queryArgs := pgx.NamedArgs{
			"scheduled_maintenance_id": param.ScheduledMaintenanceID,
			"component_id":             param.ComponentID,
		}

		preparedQuery := batchInserts.Queue(insertScheduleMaintenanceComponentsQuery, queryArgs)

		preparedQuery.Query(func(rows pgx.Rows) error {
			component, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByNameLax[scheduleMaintenanceComponentDto])
			if err != nil {
				r.lg.ErrorContext(ctx, "error collecting schedule maintenance component from batch", slog.Any("err", err))
				return err
			}

			componentsResponse = append(componentsResponse, *component)
			return nil
		})
	}

	err := r.writeDB.SendBatch(ctx, batchInserts).Close()
	if err != nil {
		r.lg.ErrorContext(ctx, "error while bulk inserting schedule maintenance components", slog.Any("err", err))
		return nil, err
	}

	result := lo.Map(componentsResponse, func(item scheduleMaintenanceComponentDto, _ int) scheduledmaintenance.ScheduledMaintenanceComponentResult {
		return scheduledmaintenance.ScheduledMaintenanceComponentResult{
			ID:                     item.ID,
			ScheduledMaintenanceID: item.ScheduledMaintenanceID,
			ComponentID:            item.ComponentID,
			CreatedAt:              item.CreatedAt,
			UpdatedAt:              item.UpdatedAt,
			DeletedAt:              nullable.SetValue(item.DeletedAt.Time, item.DeletedAt.Valid),
		}
	})

	return result, nil
}
