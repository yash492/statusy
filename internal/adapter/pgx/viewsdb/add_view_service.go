package viewsdb

import (
	"context"
	_ "embed"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/internal/common/apperrors"
	"github.com/yash492/statusy/internal/domain/views"
)

//go:embed queries/insert_views_services.sql
var insertViewServiceQuery string

//go:embed queries/insert_views_services_components.sql
var insertViewServiceComponentQuery string

//go:embed queries/insert_views_services_component_groups.sql
var insertViewServiceComponentGroupQuery string

func (r *PostgresViewsRepository) AddViewService(ctx context.Context, vs views.ViewService, componentIDs []int, componentGroupIDs []int) (views.ViewService, error) {
	tx, err := r.writeDB.Begin(ctx)
	if err != nil {
		r.lg.ErrorContext(ctx, "error starting transaction for add view service", slog.Any("err", err))
		return views.ViewService{}, apperrors.InternalError("failed to start add view service transaction", err)
	}
	defer tx.Rollback(ctx)

	// Insert the view_services row
	rows, err := tx.Query(ctx, insertViewServiceQuery, pgx.NamedArgs{
		"view_id":                vs.ViewID,
		"service_id":             vs.ServiceID,
		"include_all_components": vs.IncludeAllComponents,
	})
	if err != nil {
		r.lg.ErrorContext(ctx, "error inserting view service", slog.Uint64("view_id", uint64(vs.ViewID)), slog.Uint64("service_id", uint64(vs.ServiceID)), slog.Any("err", err))
		return views.ViewService{}, apperrors.InternalError("failed to insert view service", err)
	}
	defer rows.Close()

	dto, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[viewServiceFullDto])
	if err != nil {
		r.lg.ErrorContext(ctx, "error collecting inserted view service row", slog.Any("err", err))
		return views.ViewService{}, apperrors.InternalError("failed to collect inserted view service row", err)
	}

	// Insert component selections if not include_all_components
	if !vs.IncludeAllComponents {
		for _, componentID := range componentIDs {
			_, err := tx.Exec(ctx, insertViewServiceComponentQuery, pgx.NamedArgs{
				"view_service_id": dto.ID,
				"component_id":    componentID,
			})
			if err != nil {
				r.lg.ErrorContext(ctx, "error inserting view service component", slog.Uint64("view_service_id", uint64(dto.ID)), slog.Int("component_id", componentID), slog.Any("err", err))
				return views.ViewService{}, apperrors.InternalError("failed to insert view service component", err)
			}
		}

		for _, componentGroupID := range componentGroupIDs {
			_, err := tx.Exec(ctx, insertViewServiceComponentGroupQuery, pgx.NamedArgs{
				"view_service_id":    dto.ID,
				"component_group_id": componentGroupID,
			})
			if err != nil {
				r.lg.ErrorContext(ctx, "error inserting view service component group", slog.Uint64("view_service_id", uint64(dto.ID)), slog.Int("component_group_id", componentGroupID), slog.Any("err", err))
				return views.ViewService{}, apperrors.InternalError("failed to insert view service component group", err)
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		r.lg.ErrorContext(ctx, "error committing add view service transaction", slog.Any("err", err))
		return views.ViewService{}, apperrors.InternalError("failed to commit add view service transaction", err)
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
