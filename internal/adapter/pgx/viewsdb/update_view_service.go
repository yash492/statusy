package viewsdb

import (
	"context"
	_ "embed"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/internal/domain/views"
)

//go:embed queries/update_view_service.sql
var updateViewServiceQuery string

//go:embed queries/hard_delete_view_service_components.sql
var hardDeleteViewServiceComponentsQuery string

//go:embed queries/hard_delete_view_service_component_groups.sql
var hardDeleteViewServiceComponentGroupsQuery string

func (r *PostgresViewsRepository) UpdateViewService(ctx context.Context, vs views.ViewService, componentIDs []int, componentGroupIDs []int) (views.ViewService, error) {
	tx, err := r.writeDB.Begin(ctx)
	if err != nil {
		r.lg.ErrorContext(ctx, "error starting transaction for update view service", slog.Any("err", err))
		return views.ViewService{}, err
	}
	defer tx.Rollback(ctx)

	// Update the view_services row
	rows, err := tx.Query(ctx, updateViewServiceQuery, pgx.NamedArgs{
		"id":                     vs.ID,
		"include_all_components": vs.IncludeAllComponents,
	})
	if err != nil {
		r.lg.ErrorContext(ctx, "error updating view service", slog.Uint64("id", uint64(vs.ID)), slog.Any("err", err))
		return views.ViewService{}, err
	}
	defer rows.Close()

	dto, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[viewServiceFullDto])
	if err != nil {
		r.lg.ErrorContext(ctx, "error collecting updated view service row", slog.Any("err", err))
		return views.ViewService{}, err
	}

	// Hard-delete old component/group selections
	_, err = tx.Exec(ctx, hardDeleteViewServiceComponentsQuery, pgx.NamedArgs{
		"view_service_id": dto.ID,
	})
	if err != nil {
		r.lg.ErrorContext(ctx, "error hard-deleting view service components", slog.Uint64("view_service_id", uint64(dto.ID)), slog.Any("err", err))
		return views.ViewService{}, err
	}

	_, err = tx.Exec(ctx, hardDeleteViewServiceComponentGroupsQuery, pgx.NamedArgs{
		"view_service_id": dto.ID,
	})
	if err != nil {
		r.lg.ErrorContext(ctx, "error hard-deleting view service component groups", slog.Uint64("view_service_id", uint64(dto.ID)), slog.Any("err", err))
		return views.ViewService{}, err
	}

	// Re-insert component selections if not include_all_components
	if !vs.IncludeAllComponents {
		for _, componentID := range componentIDs {
			_, err := tx.Exec(ctx, insertViewServiceComponentQuery, pgx.NamedArgs{
				"view_service_id": dto.ID,
				"component_id":    componentID,
			})
			if err != nil {
				r.lg.ErrorContext(ctx, "error inserting view service component", slog.Uint64("view_service_id", uint64(dto.ID)), slog.Int("component_id", componentID), slog.Any("err", err))
				return views.ViewService{}, err
			}
		}

		for _, componentGroupID := range componentGroupIDs {
			_, err := tx.Exec(ctx, insertViewServiceComponentGroupQuery, pgx.NamedArgs{
				"view_service_id":  dto.ID,
				"component_group_id": componentGroupID,
			})
			if err != nil {
				r.lg.ErrorContext(ctx, "error inserting view service component group", slog.Uint64("view_service_id", uint64(dto.ID)), slog.Int("component_group_id", componentGroupID), slog.Any("err", err))
				return views.ViewService{}, err
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		r.lg.ErrorContext(ctx, "error committing update view service transaction", slog.Any("err", err))
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
