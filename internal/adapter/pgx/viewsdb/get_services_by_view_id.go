package viewsdb

import (
	"context"
	_ "embed"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/samber/lo"
	"github.com/yash492/statusy/internal/common/apperrors"
	"github.com/yash492/statusy/internal/domain/views"
)

//go:embed queries/get_view_services.sql
var getViewServicesQuery string

func (r *PostgresViewsRepository) GetServicesByViewID(ctx context.Context, viewID uint) ([]views.ViewServiceStatus, error) {
	rows, err := r.readDB.Query(ctx, getViewServicesQuery, pgx.NamedArgs{"view_id": viewID})
	if err != nil {
		r.lg.ErrorContext(ctx, "error querying view services", slog.Uint64("view_id", uint64(viewID)), slog.Any("err", err))
		return nil, apperrors.InternalError("failed to query view services", err)
	}
	defer rows.Close()

	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[viewServiceDto])
	if err != nil {
		r.lg.ErrorContext(ctx, "error collecting view services rows", slog.Uint64("view_id", uint64(viewID)), slog.Any("err", err))
		return nil, apperrors.InternalError("failed to collect view services rows", err)
	}

	return lo.Map(dtos, func(dto viewServiceDto, _ int) views.ViewServiceStatus {
		return views.ViewServiceStatus{
			ID:                   dto.ID,
			Name:                 dto.Name,
			Slug:                 dto.Slug,
			Status:               "up",
			LastIncident:         "",
			IncludeAllComponents: dto.IncludeAllComponents,
		}
	}), nil
}
