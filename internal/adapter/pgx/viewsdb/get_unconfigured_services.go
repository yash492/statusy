package viewsdb

import (
	"context"
	_ "embed"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/internal/common/apperrors"
	"github.com/yash492/statusy/internal/domain/services"
)

//go:embed queries/get_unconfigured_services.sql
var getUnconfiguredServicesQuery string

func (r *PostgresViewsRepository) GetUnconfiguredServices(ctx context.Context, viewID uint, search string) ([]services.ServiceResult, error) {
	rows, err := r.readDB.Query(ctx, getUnconfiguredServicesQuery, pgx.NamedArgs{
		"view_id": viewID,
		"search":  search,
	})
	if err != nil {
		r.lg.ErrorContext(ctx, "error querying unconfigured services", slog.Uint64("view_id", uint64(viewID)), slog.Any("err", err))
		return nil, apperrors.InternalError("failed to query unconfigured services", err)
	}
	defer rows.Close()

	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[services.ServiceResult])
	if err != nil {
		r.lg.ErrorContext(ctx, "error collecting unconfigured services rows", slog.Uint64("view_id", uint64(viewID)), slog.Any("err", err))
		return nil, apperrors.InternalError("failed to collect unconfigured services rows", err)
	}

	return dtos, nil
}
