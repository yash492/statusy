package componentsdb

import (
	"context"
	_ "embed"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/samber/lo"
	"github.com/yash492/statusy/internal/common/apperrors"
	"github.com/yash492/statusy/internal/common/nullable"
	"github.com/yash492/statusy/internal/domain/components"
)

//go:embed queries/insert_component.sql
var insertComponentQuery string

type componentDto struct {
	ID               uint
	Name             string
	ProviderID       string
	ServiceID        uint
	ComponentGroupID pgtype.Uint64
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        pgtype.Timestamptz
}

func (c *PostgresComponentRepository) SaveAll(ctx context.Context, params []components.ComponentParams) ([]components.ComponentResult, error) {
	batchInserts := &pgx.Batch{}
	componentsResponse := []componentDto{}

	for _, component := range params {
		cgVal, cgOk := component.ComponentGroupID.Get()

		queryArgs := pgx.NamedArgs{
			"name":        component.Name,
			"provider_id": component.ProviderID,
			"service_id":  component.ServiceID,
			"component_group_id": pgtype.Uint64{
				Uint64: uint64(cgVal),
				Valid:  cgOk,
			},
		}

		preparedQuery := batchInserts.Queue(
			insertComponentQuery,
			queryArgs,
		)

		preparedQuery.Query(func(rows pgx.Rows) error {
			componentRow, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[componentDto])
			if err != nil {
				c.lg.ErrorContext(ctx, "error collecting component from batch", slog.String("provider_id", component.ProviderID), slog.Uint64("service_id", uint64(component.ServiceID)), slog.Any("err", err))
				return apperrors.InternalError("failed to collect component from batch", err)
			}

			componentsResponse = append(componentsResponse, componentRow)
			return nil
		})

	}

	err := c.writeDB.SendBatch(ctx, batchInserts).Close()
	if err != nil {
		c.lg.ErrorContext(ctx, "error while bulk inserting services", slog.Any("err", err))
		return nil, apperrors.InternalError("failed to bulk insert components", err)
	}

	response := lo.Map(componentsResponse, func(item componentDto, _ int) components.ComponentResult {
		return components.ComponentResult{
			ID:               item.ID,
			Name:             item.Name,
			ProviderID:       item.ProviderID,
			ServiceID:        item.ServiceID,
			ComponentGroupID: nullable.SetValue(uint(item.ComponentGroupID.Uint64), item.ComponentGroupID.Valid),
			CreatedAt:        item.CreatedAt,
			UpdatedAt:        item.UpdatedAt,
			DeletedAt:        nullable.SetValue(item.DeletedAt.Time, item.DeletedAt.Valid),
		}
	})

	return response, nil
}
