package viewsdb

import (
	"context"
	_ "embed"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/internal/common/apperrors"
	"github.com/yash492/statusy/internal/domain/views"
)

//go:embed queries/update_view.sql
var updateViewQuery string

//go:embed queries/unset_other_default_views.sql
var unsetOtherDefaultViewsQuery string

func (r *PostgresViewsRepository) UpdateView(ctx context.Context, view views.View) (views.View, error) {
	tx, err := r.writeDB.Begin(ctx)
	if err != nil {
		r.lg.ErrorContext(ctx, "error starting transaction for update view", slog.Any("err", err))
		return views.View{}, apperrors.InternalError("failed to start update view transaction", err)
	}
	defer tx.Rollback(ctx)

	if view.IsDefault {
		_, err = tx.Exec(ctx, unsetOtherDefaultViewsQuery, pgx.NamedArgs{
			"id": view.ID,
		})
		if err != nil {
			r.lg.ErrorContext(ctx, "error unsetting other default views", slog.Uint64("id", uint64(view.ID)), slog.Any("err", err))
			return views.View{}, apperrors.InternalError("failed to unset other default views", err)
		}
	}

	rows, err := tx.Query(ctx, updateViewQuery, pgx.NamedArgs{
		"id":          view.ID,
		"name":        view.Name,
		"slug":        view.Slug,
		"description": view.Description,
		"is_default":  view.IsDefault,
	})
	if err != nil {
		r.lg.ErrorContext(ctx, "error updating view", slog.Uint64("id", uint64(view.ID)), slog.Any("err", err))
		return views.View{}, apperrors.InternalError("failed to update view", err)
	}
	defer rows.Close()

	dto, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[viewDto])
	if err != nil {
		r.lg.ErrorContext(ctx, "error collecting updated view row", slog.Uint64("id", uint64(view.ID)), slog.Any("err", err))
		return views.View{}, apperrors.InternalError("failed to collect updated view row", err)
	}

	if err := tx.Commit(ctx); err != nil {
		r.lg.ErrorContext(ctx, "error committing update view transaction", slog.Any("err", err))
		return views.View{}, apperrors.InternalError("failed to commit update view transaction", err)
	}

	servicesList, err := r.GetServicesByViewID(ctx, dto.ID)
	if err != nil {
		return views.View{}, err
	}

	return views.View{
		ID:          dto.ID,
		Name:        dto.Name,
		Slug:        dto.Slug,
		Description: dto.Description,
		IsDefault:   dto.IsDefault,
		Services:    servicesList,
		CreatedAt:   dto.CreatedAt,
		UpdatedAt:   dto.UpdatedAt,
	}, nil
}
