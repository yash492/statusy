package components

import (
	"context"
	_ "embed"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/internal/repository/components"
)

//go:embed queries/insert_component.sql
var insertComponentQuery string

func (c *PostgresComponentRepository) SaveAll(ctx context.Context, params []components.ComponentParams) ([]components.ComponentResult, error) {
	batchInserts := &pgx.Batch{}
	componentsResponse := []components.ComponentResult{}

	for _, component := range params {
		queryArgs := pgx.NamedArgs{
			"name":               component.Name,
			"provider_id":        component.ProviderID,
			"service_id":         component.ServiceID,
			"component_group_id": component.ComponentGroupID,
		}

		preparedQuery := batchInserts.Queue(
			insertComponentQuery,
			queryArgs,
		)

		preparedQuery.Query(func(rows pgx.Rows) error {
			component, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByNameLax[components.ComponentResult])
			if err != nil {
				c.lg.ErrorContext(ctx, "error collecting service %s from batch", component.Name, slog.Any("err", err))
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
	return componentsResponse, nil
}
