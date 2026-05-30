package viewsdb

import (
	"context"
	_ "embed"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/internal/domain/services"
	"github.com/yash492/statusy/internal/domain/views"
)

//go:embed queries/get_default_view.sql
var getDefaultViewQuery string

//go:embed queries/get_view_services.sql
var getViewServicesQuery string

//go:embed queries/insert_views.sql
var insertViewQuery string

//go:embed queries/get_view_by_slug.sql
var getViewBySlugQuery string

//go:embed queries/get_unconfigured_services.sql
var getUnconfiguredServicesQuery string

func (r *PostgresViewsRepository) GetDefault(ctx context.Context) (views.View, error) {
	rows, err := r.readDB.Query(ctx, getDefaultViewQuery)
	if err != nil {
		r.lg.ErrorContext(ctx, "error querying default view", slog.Any("err", err))
		return views.View{}, err
	}
	defer rows.Close()

	dto, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[viewDto])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return views.View{}, pgx.ErrNoRows
		}
		r.lg.ErrorContext(ctx, "error collecting default view row", slog.Any("err", err))
		return views.View{}, err
	}

	services, err := r.GetServicesByViewID(ctx, dto.ID)
	if err != nil {
		return views.View{}, err
	}

	return views.View{
		ID:          dto.ID,
		Name:        dto.Name,
		Slug:        dto.Slug,
		Description: dto.Description,
		IsDefault:   dto.IsDefault,
		Services:    services,
		CreatedAt:   dto.CreatedAt,
		UpdatedAt:   dto.UpdatedAt,
	}, nil
}

func (r *PostgresViewsRepository) Save(ctx context.Context, view views.View) (views.View, error) {
	rows, err := r.writeDB.Query(ctx, insertViewQuery, pgx.NamedArgs{
		"name":        view.Name,
		"slug":        view.Slug,
		"description": view.Description,
		"is_default":  view.IsDefault,
	})
	if err != nil {
		r.lg.ErrorContext(ctx, "error inserting view", slog.String("slug", view.Slug), slog.Any("err", err))
		return views.View{}, err
	}
	defer rows.Close()

	dto, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[viewDto])
	if err != nil {
		r.lg.ErrorContext(ctx, "error collecting inserted view row", slog.String("slug", view.Slug), slog.Any("err", err))
		return views.View{}, err
	}

	return views.View{
		ID:          dto.ID,
		Name:        dto.Name,
		Slug:        dto.Slug,
		Description: dto.Description,
		IsDefault:   dto.IsDefault,
		Services:    []views.ViewServiceStatus{},
		CreatedAt:   dto.CreatedAt,
		UpdatedAt:   dto.UpdatedAt,
	}, nil
}

func (r *PostgresViewsRepository) GetServicesByViewID(ctx context.Context, viewID uint) ([]views.ViewServiceStatus, error) {
	rows, err := r.readDB.Query(ctx, getViewServicesQuery, pgx.NamedArgs{"view_id": viewID})
	if err != nil {
		r.lg.ErrorContext(ctx, "error querying view services", slog.Uint64("view_id", uint64(viewID)), slog.Any("err", err))
		return nil, err
	}
	defer rows.Close()

	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[viewServiceDto])
	if err != nil {
		r.lg.ErrorContext(ctx, "error collecting view services rows", slog.Uint64("view_id", uint64(viewID)), slog.Any("err", err))
		return nil, err
	}

	result := make([]views.ViewServiceStatus, len(dtos))
	for i, dto := range dtos {
		result[i] = views.ViewServiceStatus{
			ID:                   dto.ID,
			Name:                 dto.Name,
			Slug:                 dto.Slug,
			Status:               "up",
			LastIncident:         "",
			IncludeAllComponents: dto.IncludeAllComponents,
		}
	}
	return result, nil
}

func (r *PostgresViewsRepository) GetBySlug(ctx context.Context, slug string) (views.View, error) {
	rows, err := r.readDB.Query(ctx, getViewBySlugQuery, pgx.NamedArgs{"slug": slug})
	if err != nil {
		r.lg.ErrorContext(ctx, "error querying view by slug", slog.String("slug", slug), slog.Any("err", err))
		return views.View{}, err
	}
	defer rows.Close()

	dto, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[viewDto])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return views.View{}, pgx.ErrNoRows
		}
		r.lg.ErrorContext(ctx, "error collecting view row by slug", slog.String("slug", slug), slog.Any("err", err))
		return views.View{}, err
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

func (r *PostgresViewsRepository) GetUnconfiguredServices(ctx context.Context, viewID uint) ([]services.ServiceResult, error) {
	rows, err := r.readDB.Query(ctx, getUnconfiguredServicesQuery, pgx.NamedArgs{"view_id": viewID})
	if err != nil {
		r.lg.ErrorContext(ctx, "error querying unconfigured services", slog.Uint64("view_id", uint64(viewID)), slog.Any("err", err))
		return nil, err
	}
	defer rows.Close()

	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[services.ServiceResult])
	if err != nil {
		r.lg.ErrorContext(ctx, "error collecting unconfigured services rows", slog.Uint64("view_id", uint64(viewID)), slog.Any("err", err))
		return nil, err
	}

	return dtos, nil
}

