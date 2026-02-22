package componentgroups

import (
	"context"
	_ "embed"
	"log/slog"

	"github.com/jackc/pgx/v5"
	domaincomponents "github.com/yash492/statusy/internal/domain/components"
)

//go:embed queries/insert_component_groups.sql
var insertComponentsGroupQuery string

func (c *PostgresComponentGroupsRepository) SaveAll(ctx context.Context, params []domaincomponents.GroupParams) ([]domaincomponents.GroupResult, error) {
	batchInserts := &pgx.Batch{}
	componentsGroupResponse := []domaincomponents.GroupResult{}

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
			componentGroup, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByNameLax[domaincomponents.GroupResult])
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
