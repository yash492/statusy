package componentgroups

import (
	"context"
	_ "embed"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/internal/repository/componentgroups"
)

//go:embed queries/insert_component_groups.sql
var insertComponentsGroupQuery string

func (c *PostgresComponentGroupsRepository) SaveAll(ctx context.Context, params []componentgroups.ComponentsGroupParams) ([]componentgroups.ComponentsGroupResult, error) {
	batchInserts := &pgx.Batch{}
	componentsGroupResponse := []componentgroups.ComponentsGroupResult{}

	for _, component := range params {
		queryArgs := pgx.NamedArgs{
			"name":        component.Name,
			"provider_id": component.ProviderID,
			"service_id":  component.ServiceID,
		}

		preparedQuery := batchInserts.Queue(
			insertComponentsGroupQuery,
			queryArgs,
		)

		preparedQuery.Query(func(rows pgx.Rows) error {
			componentGroup, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByNameLax[componentgroups.ComponentsGroupResult])
			if err != nil {
				c.lg.ErrorContext(ctx, "error collecting service %s from batch", componentGroup.Name, slog.Any("err", err))
				return err
			}

			componentsGroupResponse = append(componentsGroupResponse, *componentGroup)
			return nil
		})

	}

	err := c.writeDB.SendBatch(ctx, batchInserts).Close()
	if err != nil {
		c.lg.ErrorContext(ctx, "error while bulk inserting services", slog.Any("err", err))
		return nil, err
	}
	return componentsGroupResponse, nil
}
