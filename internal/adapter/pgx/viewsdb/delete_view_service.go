package viewsdb

import (
	"context"
	_ "embed"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/internal/common/apperrors"
)

//go:embed queries/soft_delete_view_service_components.sql
var softDeleteViewServiceComponentsQuery string

//go:embed queries/soft_delete_view_service_component_groups.sql
var softDeleteViewServiceComponentGroupsQuery string

//go:embed queries/soft_delete_view_service.sql
var softDeleteViewServiceQuery string

func (r *PostgresViewsRepository) DeleteViewService(ctx context.Context, viewID uint, serviceID uint) error {
	tx, err := r.writeDB.Begin(ctx)
	if err != nil {
		r.lg.ErrorContext(ctx, "error starting transaction for delete view service", slog.Any("err", err))
		return apperrors.InternalError("failed to start delete view service transaction", err)
	}
	defer tx.Rollback(ctx)

	// Get view service to obtain the view_service_id
	rows, err := tx.Query(ctx, getViewServiceQuery, pgx.NamedArgs{
		"view_id":    viewID,
		"service_id": serviceID,
	})
	if err != nil {
		r.lg.ErrorContext(ctx, "error querying view service for delete", slog.Uint64("view_id", uint64(viewID)), slog.Uint64("service_id", uint64(serviceID)), slog.Any("err", err))
		return apperrors.InternalError("failed to query view service for delete", err)
	}
	defer rows.Close()

	dto, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[viewServiceFullDto])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apperrors.NotFoundError("view service not found", err)
		}
		r.lg.ErrorContext(ctx, "error collecting view service row for delete", slog.Any("err", err))
		return apperrors.InternalError("failed to collect view service row for delete", err)
	}

	// Soft-delete component rows
	_, err = tx.Exec(ctx, softDeleteViewServiceComponentsQuery, pgx.NamedArgs{
		"view_service_id": dto.ID,
	})
	if err != nil {
		r.lg.ErrorContext(ctx, "error soft-deleting view service components", slog.Uint64("view_service_id", uint64(dto.ID)), slog.Any("err", err))
		return apperrors.InternalError("failed to soft-delete view service components", err)
	}

	// Soft-delete component group rows
	_, err = tx.Exec(ctx, softDeleteViewServiceComponentGroupsQuery, pgx.NamedArgs{
		"view_service_id": dto.ID,
	})
	if err != nil {
		r.lg.ErrorContext(ctx, "error soft-deleting view service component groups", slog.Uint64("view_service_id", uint64(dto.ID)), slog.Any("err", err))
		return apperrors.InternalError("failed to soft-delete view service component groups", err)
	}

	// Soft-delete the view_services row
	_, err = tx.Exec(ctx, softDeleteViewServiceQuery, pgx.NamedArgs{
		"view_id":    viewID,
		"service_id": serviceID,
	})
	if err != nil {
		r.lg.ErrorContext(ctx, "error soft-deleting view service", slog.Uint64("view_id", uint64(viewID)), slog.Uint64("service_id", uint64(serviceID)), slog.Any("err", err))
		return apperrors.InternalError("failed to soft-delete view service", err)
	}

	if err := tx.Commit(ctx); err != nil {
		r.lg.ErrorContext(ctx, "error committing delete view service transaction", slog.Any("err", err))
		return apperrors.InternalError("failed to commit delete view service transaction", err)
	}

	return nil
}
