package componentsdb

import (
	"context"
	_ "embed"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/samber/lo"
	"github.com/yash492/statusy/internal/common/nullable"
	"github.com/yash492/statusy/internal/domain/components"
)

//go:embed queries/get_components_by_service.sql
var getComponentsByServiceQuery string

func (r *PostgresComponentRepository) GetByServiceID(ctx context.Context, serviceID uint) ([]components.ComponentResult, error) {
	rows, err := r.readDB.Query(ctx, getComponentsByServiceQuery, pgx.NamedArgs{"service_id": serviceID})
	if err != nil {
		r.lg.ErrorContext(ctx, "error querying components by service id", slog.Uint64("service_id", uint64(serviceID)), slog.Any("err", err))
		return nil, err
	}
	defer rows.Close()

	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[componentDto])
	if err != nil {
		r.lg.ErrorContext(ctx, "error collecting components rows", slog.Uint64("service_id", uint64(serviceID)), slog.Any("err", err))
		return nil, err
	}

	result := lo.Map(dtos, func(item componentDto, _ int) components.ComponentResult {
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

	return result, nil
}
