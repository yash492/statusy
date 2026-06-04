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

	rows, err := tx.Query(ctx, insertViewServiceQuery, pgx.NamedArgs{
		"view_id":                        vs.ViewID,
		"service_id":                     vs.ServiceID,
		"include_all_components":         vs.IncludeAllComponents,
		"monitor_incidents":              vs.MonitorIncidents,
		"monitor_scheduled_maintenances": vs.MonitorScheduledMaintenances,
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
	if !vs.IncludeAllComponents && (len(componentIDs) > 0 || len(componentGroupIDs) > 0) {
		batchInserts := &pgx.Batch{}

		for _, componentID := range componentIDs {
			batchInserts.Queue(insertViewServiceComponentQuery, pgx.NamedArgs{
				"view_service_id": dto.ID,
				"component_id":    componentID,
			})
		}

		for _, componentGroupID := range componentGroupIDs {
			batchInserts.Queue(insertViewServiceComponentGroupQuery, pgx.NamedArgs{
				"view_service_id":    dto.ID,
				"component_group_id": componentGroupID,
			})
		}

		batchResults := tx.SendBatch(ctx, batchInserts)
		if err := batchResults.Close(); err != nil {
			r.lg.ErrorContext(ctx, "error executing batch inserts for view service components/groups", slog.Uint64("view_service_id", uint64(dto.ID)), slog.Any("err", err))
			return views.ViewService{}, apperrors.InternalError("failed to execute batch inserts for components/groups", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		r.lg.ErrorContext(ctx, "error committing add view service transaction", slog.Any("err", err))
		return views.ViewService{}, apperrors.InternalError("failed to commit add view service transaction", err)
	}

	compIDs := componentIDs
	if compIDs == nil {
		compIDs = []int{}
	}

	compGrpIDs := componentGroupIDs
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
