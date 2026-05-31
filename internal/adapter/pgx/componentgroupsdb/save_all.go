package componentgroupsdb

import (
	"context"
	_ "embed"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/samber/lo"
	"github.com/yash492/statusy/internal/common/apperrors"
	"github.com/yash492/statusy/internal/common/nullable"
	"github.com/yash492/statusy/internal/domain/components"
)

//go:embed queries/insert_component_groups.sql
var insertComponentgroupQuery string

func (c *PostgresComponentGroupsRepository) SaveAll(ctx context.Context, params []components.GroupParams) ([]components.ComponentGroupResult, error) {
	batchInserts := &pgx.Batch{}
	componentgroupResponse := []componentGroupDto{}

	for _, component := range params {
		queryArgs := pgx.NamedArgs{
			"name":        component.Name,
			"provider_id": component.ProviderID,
			"service_id":  component.ServiceID,
		}

		preparedQuery := batchInserts.Queue(
			insertComponentgroupQuery,
			queryArgs,
		)

		preparedQuery.Query(func(rows pgx.Rows) error {
			componentGroupRow, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[componentGroupDto])
			if err != nil {
				c.lg.ErrorContext(ctx, "error collecting component group from batch", slog.String("name", component.Name), slog.Any("err", err))
				return apperrors.InternalError("failed to collect component group from batch", err)
			}

			componentgroupResponse = append(componentgroupResponse, componentGroupRow)
			return nil
		})

	}

	err := c.writeDB.SendBatch(ctx, batchInserts).Close()
	if err != nil {
		c.lg.ErrorContext(ctx, "error while bulk inserting services", slog.Any("err", err))
		return nil, apperrors.InternalError("failed to bulk insert component groups", err)
	}

	response := lo.Map(componentgroupResponse, func(item componentGroupDto, _ int) components.ComponentGroupResult {
		return components.ComponentGroupResult{
			ID:         item.ID,
			Name:       item.Name,
			ProviderID: item.ProviderID,
			ServiceID:  item.ServiceID,
			CreatedAt:  item.CreatedAt,
			UpdatedAt:  item.CreatedAt,
			DeletedAt:  nullable.SetValue(item.DeletedAt.Time, item.DeletedAt.Valid),
		}
	})
	return response, nil
}
