package viewsdb

import (
	"context"
	_ "embed"
	"errors"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/internal/common/apperrors"
	"github.com/yash492/statusy/internal/domain/views"
)

//go:embed queries/get_view_service.sql
var getViewServiceQuery string

// viewServiceFullDto maps a full row from the view_services table.
type viewServiceFullDto struct {
	ID                           uint       `db:"id"`
	ViewID                       uint       `db:"view_id"`
	ServiceID                    uint       `db:"service_id"`
	IncludeAllComponents         bool       `db:"include_all_components"`
	MonitorIncidents             bool       `db:"monitor_incidents"`
	MonitorScheduledMaintenances bool       `db:"monitor_scheduled_maintenances"`
	ComponentIDs                 []int      `db:"component_ids"`
	ComponentGroupIDs            []int      `db:"component_group_ids"`
	CreatedAt                    time.Time  `db:"created_at"`
	UpdatedAt                    time.Time  `db:"updated_at"`
	DeletedAt                    *time.Time `db:"deleted_at"`
}
func (r *PostgresViewsRepository) GetViewService(ctx context.Context, viewID uint, serviceID uint) (views.ViewService, error) {
	rows, err := r.readDB.Query(ctx, getViewServiceQuery, pgx.NamedArgs{
		"view_id":    viewID,
		"service_id": serviceID,
	})
	if err != nil {
		r.lg.ErrorContext(ctx, "error querying view service", slog.Uint64("view_id", uint64(viewID)), slog.Uint64("service_id", uint64(serviceID)), slog.Any("err", err))
		return views.ViewService{}, apperrors.InternalError("failed to query view service", err)
	}
	defer rows.Close()

	dto, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[viewServiceFullDto])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return views.ViewService{}, apperrors.NotFoundError("view service not found", err)
		}
		r.lg.ErrorContext(ctx, "error collecting view service row", slog.Uint64("view_id", uint64(viewID)), slog.Uint64("service_id", uint64(serviceID)), slog.Any("err", err))
		return views.ViewService{}, apperrors.InternalError("failed to collect view service row", err)
	}

	compIDs := dto.ComponentIDs
	if compIDs == nil {
		compIDs = []int{}
	}

	compGrpIDs := dto.ComponentGroupIDs
	if compGrpIDs == nil {
		compGrpIDs = []int{}
	}

	return views.ViewService{
		ID:                           dto.ID,
		ViewID:                       dto.ViewID,
		ServiceID:                    dto.ServiceID,
		IncludeAllComponents:         dto.IncludeAllComponents,
		MonitorIncidents:             dto.MonitorIncidents,
		MonitorScheduledMaintenances: dto.MonitorScheduledMaintenances,
		ComponentIDs:                 compIDs,
		ComponentGroupIDs:            compGrpIDs,
		CreatedAt:                    dto.CreatedAt,
		UpdatedAt:                    dto.UpdatedAt,
	}, nil
}
