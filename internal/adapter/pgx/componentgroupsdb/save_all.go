package componentgroupsdb

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

//go:embed queries/insert_component_groups.sql
var insertComponentgroupQuery string

type componentGroupDto struct {
	ID         uint             `db:"id"`
	Name       string           `db:"name"`
	ProviderID string           `db:"provider_id"`
	ServiceID  uint             `db:"service_id"`
	CreatedAt  time.Time        `db:"created_at"`
	UpdatedAt  time.Time        `db:"updated_at"`
	DeletedAt  pgtype.Timestamp `db:"deleted_at"`
}

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
			componentGroup, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByNameLax[componentGroupDto])
			if err != nil {
				c.lg.ErrorContext(ctx, "error collecting service %s from batch", componentGroup.Name, slog.Any("err", err))
				return err
			}

			componentgroupResponse = append(componentgroupResponse, *componentGroup)
			return nil
		})

	}

	err := c.writeDB.SendBatch(ctx, batchInserts).Close()
	if err != nil {
		c.lg.ErrorContext(ctx, "error while bulk inserting services", slog.Any("err", err))
		return nil, err
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
