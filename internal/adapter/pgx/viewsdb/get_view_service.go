package viewsdb

import (
	"context"
	_ "embed"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/internal/domain/views"
)

//go:embed queries/get_view_service.sql
var getViewServiceQuery string

func (r *PostgresViewsRepository) GetViewService(ctx context.Context, viewID uint, serviceID uint) (views.ViewService, error) {
	rows, err := r.readDB.Query(ctx, getViewServiceQuery, pgx.NamedArgs{
		"view_id":    viewID,
		"service_id": serviceID,
	})
	if err != nil {
		r.lg.ErrorContext(ctx, "error querying view service", slog.Uint64("view_id", uint64(viewID)), slog.Uint64("service_id", uint64(serviceID)), slog.Any("err", err))
		return views.ViewService{}, err
	}
	defer rows.Close()

	dto, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[viewServiceFullDto])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return views.ViewService{}, pgx.ErrNoRows
		}
		r.lg.ErrorContext(ctx, "error collecting view service row", slog.Uint64("view_id", uint64(viewID)), slog.Uint64("service_id", uint64(serviceID)), slog.Any("err", err))
		return views.ViewService{}, err
	}

	return views.ViewService{
		ID:                   dto.ID,
		ViewID:               dto.ViewID,
		ServiceID:            dto.ServiceID,
		IncludeAllComponents: dto.IncludeAllComponents,
		CreatedAt:            dto.CreatedAt,
		UpdatedAt:            dto.UpdatedAt,
	}, nil
}
