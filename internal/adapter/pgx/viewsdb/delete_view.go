package viewsdb

import (
	"context"
	_ "embed"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/internal/common/apperrors"
)

//go:embed queries/hard_delete_view_service_components_by_view_id.sql
var hardDeleteViewServiceComponentsByViewIDQuery string

//go:embed queries/hard_delete_view_service_component_groups_by_view_id.sql
var hardDeleteViewServiceComponentGroupsByViewIDQuery string

//go:embed queries/soft_delete_view_services_by_view_id.sql
var softDeleteViewServicesByViewIDQuery string

//go:embed queries/soft_delete_view.sql
var softDeleteViewQuery string

func (r *PostgresViewsRepository) DeleteView(ctx context.Context, viewID uint) error {
	tx, err := r.writeDB.Begin(ctx)
	if err != nil {
		r.lg.ErrorContext(ctx, "error starting transaction for delete view", slog.Any("err", err))
		return apperrors.InternalError("failed to start delete view transaction", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, hardDeleteViewServiceComponentsByViewIDQuery, pgx.NamedArgs{
		"view_id": viewID,
	})
	if err != nil {
		r.lg.ErrorContext(ctx, "error soft-deleting view service components for view", slog.Uint64("view_id", uint64(viewID)), slog.Any("err", err))
		return apperrors.InternalError("failed to soft-delete view service components", err)
	}

	_, err = tx.Exec(ctx, hardDeleteViewServiceComponentGroupsByViewIDQuery, pgx.NamedArgs{
		"view_id": viewID,
	})
	if err != nil {
		r.lg.ErrorContext(ctx, "error soft-deleting view service component groups for view", slog.Uint64("view_id", uint64(viewID)), slog.Any("err", err))
		return apperrors.InternalError("failed to soft-delete view service component groups", err)
	}

	// Soft-delete view services
	_, err = tx.Exec(ctx, softDeleteViewServicesByViewIDQuery, pgx.NamedArgs{
		"view_id": viewID,
	})
	if err != nil {
		r.lg.ErrorContext(ctx, "error soft-deleting view services for view", slog.Uint64("view_id", uint64(viewID)), slog.Any("err", err))
		return apperrors.InternalError("failed to soft-delete view services", err)
	}

	// Soft-delete view
	_, err = tx.Exec(ctx, softDeleteViewQuery, pgx.NamedArgs{
		"view_id": viewID,
	})
	if err != nil {
		r.lg.ErrorContext(ctx, "error soft-deleting view", slog.Uint64("view_id", uint64(viewID)), slog.Any("err", err))
		return apperrors.InternalError("failed to soft-delete view", err)
	}

	if err := tx.Commit(ctx); err != nil {
		r.lg.ErrorContext(ctx, "error committing delete view transaction", slog.Any("err", err))
		return apperrors.InternalError("failed to commit delete view transaction", err)
	}

	return nil
}
