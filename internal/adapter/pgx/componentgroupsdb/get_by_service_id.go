package componentgroupsdb

import (
	"context"
	_ "embed"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/samber/lo"
	"github.com/yash492/statusy/internal/common/nullable"
	"github.com/yash492/statusy/internal/domain/components"
)

//go:embed queries/get_component_groups_by_service.sql
var getComponentGroupsByServiceQuery string

func (r *PostgresComponentGroupsRepository) GetByServiceID(ctx context.Context, serviceID uint) ([]components.ComponentGroupResult, error) {
	rows, err := r.readDB.Query(ctx, getComponentGroupsByServiceQuery, pgx.NamedArgs{"service_id": serviceID})
	if err != nil {
		r.lg.ErrorContext(ctx, "error querying component groups by service id", slog.Uint64("service_id", uint64(serviceID)), slog.Any("err", err))
		return nil, err
	}
	defer rows.Close()

	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[componentGroupDto])
	if err != nil {
		r.lg.ErrorContext(ctx, "error collecting component groups rows", slog.Uint64("service_id", uint64(serviceID)), slog.Any("err", err))
		return nil, err
	}

	result := lo.Map(dtos, func(item componentGroupDto, _ int) components.ComponentGroupResult {
		return components.ComponentGroupResult{
			ID:         item.ID,
			Name:       item.Name,
			ProviderID: item.ProviderID,
			ServiceID:  item.ServiceID,
			CreatedAt:  item.CreatedAt,
			UpdatedAt:  item.UpdatedAt,
			DeletedAt:  nullable.SetValue(item.DeletedAt.Time, item.DeletedAt.Valid),
		}
	})

	return result, nil
}
