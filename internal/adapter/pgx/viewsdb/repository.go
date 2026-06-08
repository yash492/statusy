package viewsdb

import (
	"context"
	_ "embed"
	"errors"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/samber/lo"
	"github.com/yash492/statusy/internal/common/apperrors"
	"github.com/yash492/statusy/internal/common/snowflake"
	"github.com/yash492/statusy/internal/domain/services"
	"github.com/yash492/statusy/internal/domain/views"
)

// viewDto maps a row from the views table.
type viewDto struct {
	ID          uint       `db:"id"`
	Name        string     `db:"name"`
	PublicID    string     `db:"public_id"`
	Description string     `db:"description"`
	IsDefault   bool       `db:"is_default"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
}

// viewServiceDto maps a row joining view_services → services.
type viewServiceDto struct {
	ID                           uint       `db:"id"`
	Name                         string     `db:"name"`
	Slug                         string     `db:"slug"`
	IncludeAllComponents         bool       `db:"include_all_components"`
	MonitorIncidents             bool       `db:"monitor_incidents"`
	MonitorScheduledMaintenances bool       `db:"monitor_scheduled_maintenances"`
	Status                       string     `db:"status"`
	LastIncident                 string     `db:"last_incident"`
	LastIncidentLink             string     `db:"last_incident_link"`
	UpcomingMaintenance          string     `db:"upcoming_maintenance"`
	UpcomingMaintenanceLink      string     `db:"upcoming_maintenance_link"`
	ComponentIDs                 []int      `db:"component_ids"`
	ComponentGroupIDs            []int      `db:"component_group_ids"`
	UpdatedAt                    time.Time  `db:"updated_at"`
	DeletedAt                    *time.Time `db:"deleted_at"`
}

//go:embed queries/get_default_view.sql
var getDefaultViewQuery string

//go:embed queries/get_view_services.sql
var getViewServicesQuery string

//go:embed queries/insert_views.sql
var insertViewQuery string

//go:embed queries/get_view_by_public_id.sql
var getViewByPublicIDQuery string

//go:embed queries/get_unconfigured_services.sql
var getUnconfiguredServicesQuery string

func (r *PostgresViewsRepository) GetDefault(ctx context.Context) (views.View, error) {
	rows, err := r.readDB.Query(ctx, getDefaultViewQuery)
	if err != nil {
		r.lg.ErrorContext(ctx, "error querying default view", slog.Any("err", err))
		return views.View{}, apperrors.InternalError("failed to query default view", err)
	}
	defer rows.Close()

	dto, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[viewDto])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return views.View{}, apperrors.NotFoundError("default view not found", err)
		}
		r.lg.ErrorContext(ctx, "error collecting default view row", slog.Any("err", err))
		return views.View{}, apperrors.InternalError("failed to collect default view row", err)
	}

	return views.View{
		ID:          dto.ID,
		Name:        dto.Name,
		PublicID:    dto.PublicID,
		Description: dto.Description,
		IsDefault:   dto.IsDefault,
		CreatedAt:   dto.CreatedAt,
		UpdatedAt:   dto.UpdatedAt,
	}, nil
}

func (r *PostgresViewsRepository) Save(ctx context.Context, view views.View) (views.View, error) {
	if view.PublicID == "" {
		view.PublicID = snowflake.Generate()
	}

	rows, err := r.writeDB.Query(ctx, insertViewQuery, pgx.NamedArgs{
		"name":        view.Name,
		"public_id":   view.PublicID,
		"description": view.Description,
		"is_default":  view.IsDefault,
	})
	if err != nil {
		r.lg.ErrorContext(ctx, "error inserting view", slog.String("public_id", view.PublicID), slog.Any("err", err))
		return views.View{}, apperrors.InternalError("failed to insert view", err)
	}
	defer rows.Close()

	dto, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[viewDto])
	if err != nil {
		r.lg.ErrorContext(ctx, "error collecting inserted view row", slog.String("public_id", view.PublicID), slog.Any("err", err))
		return views.View{}, apperrors.InternalError("failed to collect inserted view row", err)
	}

	return views.View{
		ID:          dto.ID,
		Name:        dto.Name,
		PublicID:    dto.PublicID,
		Description: dto.Description,
		IsDefault:   dto.IsDefault,
		CreatedAt:   dto.CreatedAt,
		UpdatedAt:   dto.UpdatedAt,
	}, nil
}

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

func (r *PostgresViewsRepository) GetByPublicID(ctx context.Context, publicID string) (views.View, error) {
	rows, err := r.readDB.Query(ctx, getViewByPublicIDQuery, pgx.NamedArgs{"public_id": publicID})
	if err != nil {
		r.lg.ErrorContext(ctx, "error querying view by public_id", slog.String("public_id", publicID), slog.Any("err", err))
		return views.View{}, apperrors.InternalError("failed to query view by public_id", err)
	}
	defer rows.Close()

	dto, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[viewDto])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return views.View{}, apperrors.NotFoundError("view not found", err)
		}
		r.lg.ErrorContext(ctx, "error collecting view row by public_id", slog.String("public_id", publicID), slog.Any("err", err))
		return views.View{}, apperrors.InternalError("failed to collect view row by public_id", err)
	}

	return views.View{
		ID:          dto.ID,
		Name:        dto.Name,
		PublicID:    dto.PublicID,
		Description: dto.Description,
		IsDefault:   dto.IsDefault,
		CreatedAt:   dto.CreatedAt,
		UpdatedAt:   dto.UpdatedAt,
	}, nil
}

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
