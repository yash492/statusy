package componentsdb

import (
	"context"
	_ "embed"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/samber/lo"
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
		queryArgs := pgx.NamedArgs{
			"name":        component.Name,
			"provider_id": component.ProviderID,
			"service_id":  component.ServiceID,
			"component_group_id": pgtype.Uint64{
				Uint64: uint64(component.ComponentGroupID.Value),
				Valid:  component.ComponentGroupID.Valid,
			},
		}

		preparedQuery := batchInserts.Queue(
			insertComponentQuery,
			queryArgs,
		)

		preparedQuery.Query(func(rows pgx.Rows) error {
			component, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByNameLax[componentDto])
			if err != nil {
				c.lg.ErrorContext(ctx, "error collecting component %s for service %s from batch", component.ProviderID, component.ServiceID, slog.Any("err", err))
				return err
			}

			componentsResponse = append(componentsResponse, *component)
			return nil
		})

	}

	err := c.writeDB.SendBatch(ctx, batchInserts).Close()
	if err != nil {
		c.lg.ErrorContext(ctx, "error while bulk inserting services", slog.Any("err", err))
		return nil, err
	}

	response := lo.Map(componentsResponse, func(item componentDto, _ int) components.ComponentResult {
		return components.ComponentResult{
			ID:         item.ID,
			Name:       item.Name,
			ProviderID: item.ProviderID,
			ServiceID:  item.ServiceID,
			ComponentGroupID: nullable.Nullable[uint]{
				Value: uint(item.ComponentGroupID.Uint64),
				Valid: item.ComponentGroupID.Valid,
			},
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
			DeletedAt: nullable.Nullable[time.Time]{
				Value: item.DeletedAt.Time,
				Valid: item.DeletedAt.Valid,
			},
		}
	})

	return response, nil
}
