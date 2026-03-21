package scheduledmaintenanceupdatesdb

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

//go:embed queries/insert_scheduled_maintenance_updates.sql
var insertScheduledMaintenanceUpdatesQuery string

type scheduledMaintenanceUpdateDto struct {
	ID                     uint             `db:"id"`
	ScheduledMaintenanceID uint             `db:"scheduled_maintenance_id"`
	Description            string           `db:"description"`
	ProviderID             string           `db:"provider_id"`
	ProviderStatus         string           `db:"provider_status"`
	Status                 string           `db:"status"`
	StatusTime             time.Time        `db:"status_time"`
	CreatedAt              time.Time        `db:"created_at"`
	UpdatedAt              time.Time        `db:"updated_at"`
	DeletedAt              pgtype.Timestamp `db:"deleted_at"`
}

func (r *PostgresScheduledMaintenanceUpdatesRepository) SaveAll(ctx context.Context, params []scheduledmaintenance.ScheduledMaintenanceUpdateParams) ([]scheduledmaintenance.ScheduledMaintenanceUpdateResult, error) {
	batchInserts := &pgx.Batch{}
	updatesResponse := []scheduledMaintenanceUpdateDto{}

	for _, param := range params {
		queryArgs := pgx.NamedArgs{
			"scheduled_maintenance_id": param.ScheduledMaintenanceID,
			"description":              param.Description,
			"provider_id":              param.ProviderID,
			"provider_status":          param.ProviderStatus,
			"status":                   param.Status,
			"status_time":              param.StatusTime,
		}

		preparedQuery := batchInserts.Queue(insertScheduledMaintenanceUpdatesQuery, queryArgs)

		preparedQuery.Query(func(rows pgx.Rows) error {
			update, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByNameLax[scheduledMaintenanceUpdateDto])
			if err != nil {
				r.lg.ErrorContext(ctx, "error collecting scheduled maintenance update from batch", slog.Any("err", err))
				return err
			}

			updatesResponse = append(updatesResponse, *update)
			return nil
		})
	}

	err := r.writeDB.SendBatch(ctx, batchInserts).Close()
	if err != nil {
		r.lg.ErrorContext(ctx, "error while bulk inserting scheduled maintenance updates", slog.Any("err", err))
		return nil, err
	}

	result := lo.Map(updatesResponse, func(item scheduledMaintenanceUpdateDto, _ int) scheduledmaintenance.ScheduledMaintenanceUpdateResult {
		return scheduledmaintenance.ScheduledMaintenanceUpdateResult{
			ID:                     item.ID,
			ScheduledMaintenanceID: item.ScheduledMaintenanceID,
			Description:            item.Description,
			ProviderID:             item.ProviderID,
			ProviderStatus:         item.ProviderStatus,
			Status:                 item.Status,
			StatusTime:             item.StatusTime,
			CreatedAt:              item.CreatedAt,
			UpdatedAt:              item.UpdatedAt,
			DeletedAt:              nullable.SetValue(item.DeletedAt.Time, item.DeletedAt.Valid),
		}
	})

	return result, nil
}
