package scheduledmaintenancesdb

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

//go:embed queries/insert_scheduled_maintenances.sql
var insertScheduledMaintenancesQuery string

type scheduledMaintenanceDto struct {
	ID                uint             `db:"id"`
	Title             string           `db:"title"`
	Link              string           `db:"link"`
	ProviderImpact    pgtype.Text      `db:"provider_impact"`
	Impact            pgtype.Text      `db:"impact"`
	StartsAt          time.Time        `db:"starts_at"`
	EndsAt            time.Time        `db:"ends_at"`
	ServiceID         uint             `db:"service_id"`
	ProviderID        string           `db:"provider_id"`
	ProviderCreatedAt time.Time        `db:"provider_created_at"`
	CreatedAt         time.Time        `db:"created_at"`
	UpdatedAt         time.Time        `db:"updated_at"`
	DeletedAt         pgtype.Timestamp `db:"deleted_at"`
}

func (c *PostgresScheduledMaintenanceRepository) SaveAll(ctx context.Context, params []scheduledmaintenance.ScheduledMaintenanceParams) ([]scheduledmaintenance.ScheduledMaintenanceResult, error) {
	batchInserts := &pgx.Batch{}
	scheduledMaintenanceResponse := []scheduledMaintenanceDto{}

	for _, param := range params {
		provImpStr, provImpOk := param.ProviderImpact.Get()
		impactStr, impactOk := param.Impact.Get()

		queryArgs := pgx.NamedArgs{
			"title": param.Title,
			"link":  param.Link,
			"provider_impact": pgtype.Text{
				String: provImpStr,
				Valid:  provImpOk,
			},
			"impact": pgtype.Text{
				String: impactStr,
				Valid:  impactOk,
			},
			"starts_at":           param.StartsAt,
			"ends_at":             param.EndsAt,
			"service_id":          param.ServiceID,
			"provider_id":         param.ProviderID,
			"provider_created_at": param.ProviderCreatedAt,
		}

		preparedQuery := batchInserts.Queue(
			insertScheduledMaintenancesQuery,
			queryArgs,
		)

		preparedQuery.Query(func(rows pgx.Rows) error {
			scheduledMaintenance, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByNameLax[scheduledMaintenanceDto])
			if err != nil {
				c.lg.ErrorContext(ctx, "error collecting scheduled maintenance %v from %v from batch", param.ProviderID, param.ServiceID, slog.Any("err", err))
				return err
			}

			scheduledMaintenanceResponse = append(scheduledMaintenanceResponse, *scheduledMaintenance)
			return nil
		})

	}

	err := c.writeDB.SendBatch(ctx, batchInserts).Close()
	if err != nil {
		c.lg.ErrorContext(ctx, "error while bulk inserting scheduled maintenances", slog.Any("err", err))
		return nil, err
	}

	response := lo.Map(scheduledMaintenanceResponse, func(item scheduledMaintenanceDto, _ int) scheduledmaintenance.ScheduledMaintenanceResult {
		return scheduledmaintenance.ScheduledMaintenanceResult{
			ID:                item.ID,
			Title:             item.Title,
			Link:              item.Link,
			ProviderImpact:    nullable.SetValue(item.ProviderImpact.String, item.ProviderImpact.Valid),
			Impact:            nullable.SetValue(item.Impact.String, item.Impact.Valid),
			StartsAt:          item.StartsAt,
			EndsAt:            item.EndsAt,
			ServiceID:         item.ServiceID,
			ProviderID:        item.ProviderID,
			ProviderCreatedAt: item.ProviderCreatedAt,
			CreatedAt:         item.CreatedAt,
			UpdatedAt:         item.UpdatedAt,
			DeletedAt:         nullable.SetValue(item.DeletedAt.Time, item.DeletedAt.Valid),
		}
	})

	return response, nil
}
