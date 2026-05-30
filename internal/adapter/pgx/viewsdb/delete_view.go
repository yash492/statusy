package viewsdb

import (
	"context"
	_ "embed"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

//go:embed queries/soft_delete_view_service_components_by_view_id.sql
var softDeleteViewServiceComponentsByViewIDQuery string

//go:embed queries/soft_delete_view_service_component_groups_by_view_id.sql
var softDeleteViewServiceComponentGroupsByViewIDQuery string

//go:embed queries/soft_delete_view_services_by_view_id.sql
var softDeleteViewServicesByViewIDQuery string

//go:embed queries/soft_delete_view.sql
var softDeleteViewQuery string

func (r *PostgresViewsRepository) DeleteView(ctx context.Context, viewID uint) error {
	tx, err := r.writeDB.Begin(ctx)
	if err != nil {
		r.lg.ErrorContext(ctx, "error starting transaction for delete view", slog.Any("err", err))
		return err
	}
	defer tx.Rollback(ctx)

	// Soft-delete view service components
	_, err = tx.Exec(ctx, softDeleteViewServiceComponentsByViewIDQuery, pgx.NamedArgs{
		"view_id": viewID,
	})
	if err != nil {
		r.lg.ErrorContext(ctx, "error soft-deleting view service components for view", slog.Uint64("view_id", uint64(viewID)), slog.Any("err", err))
		return err
	}

	// Soft-delete view service component groups
	_, err = tx.Exec(ctx, softDeleteViewServiceComponentGroupsByViewIDQuery, pgx.NamedArgs{
		"view_id": viewID,
	})
	if err != nil {
		r.lg.ErrorContext(ctx, "error soft-deleting view service component groups for view", slog.Uint64("view_id", uint64(viewID)), slog.Any("err", err))
		return err
	}

	// Soft-delete view services
	_, err = tx.Exec(ctx, softDeleteViewServicesByViewIDQuery, pgx.NamedArgs{
		"view_id": viewID,
	})
	if err != nil {
		r.lg.ErrorContext(ctx, "error soft-deleting view services for view", slog.Uint64("view_id", uint64(viewID)), slog.Any("err", err))
		return err
	}

	// Soft-delete view
	_, err = tx.Exec(ctx, softDeleteViewQuery, pgx.NamedArgs{
		"view_id": viewID,
	})
	if err != nil {
		r.lg.ErrorContext(ctx, "error soft-deleting view", slog.Uint64("view_id", uint64(viewID)), slog.Any("err", err))
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		r.lg.ErrorContext(ctx, "error committing delete view transaction", slog.Any("err", err))
		return err
	}

	return nil
}
