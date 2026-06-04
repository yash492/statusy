package viewsdb

import (
	"context"
	_ "embed"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/internal/common/apperrors"
	"github.com/yash492/statusy/internal/domain/views"
)

//go:embed queries/get_view_services_paginated.sql
var getViewServicesPaginatedQuery string

//go:embed queries/count_view_services.sql
var countViewServicesQuery string

func (r *PostgresViewsRepository) GetViewServices(ctx context.Context, viewID uint, search string, limit int, offset int) ([]views.ViewServiceStatus, int64, int64, int64, error) {
	// Get paginated services
	rows, err := r.readDB.Query(ctx, getViewServicesPaginatedQuery, pgx.NamedArgs{
		"view_id": viewID,
		"search":  search,
		"limit":   limit,
		"offset":  offset,
	})
	if err != nil {
		r.lg.ErrorContext(ctx, "error querying paginated view services", slog.Any("err", err))
		return nil, 0, 0, 0, apperrors.InternalError("failed to query paginated view services", err)
	}
	defer rows.Close()

	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[viewServiceDto])
	if err != nil {
		r.lg.ErrorContext(ctx, "error collecting paginated view services rows", slog.Any("err", err))
		return nil, 0, 0, 0, apperrors.InternalError("failed to collect paginated view services rows", err)
	}

	// Get counts (total, down, up)
	var totalCount, downCount, upCount int64
	err = r.readDB.QueryRow(ctx, countViewServicesQuery, pgx.NamedArgs{
		"view_id": viewID,
		"search":  search,
	}).Scan(&totalCount, &downCount, &upCount)
	if err != nil {
		r.lg.ErrorContext(ctx, "error counting view services", slog.Any("err", err))
		return nil, 0, 0, 0, apperrors.InternalError("failed to count view services", err)
	}

	result := make([]views.ViewServiceStatus, len(dtos))
	for i, dto := range dtos {
		result[i] = views.ViewServiceStatus{
			ID:                           dto.ID,
			Name:                         dto.Name,
			Slug:                         dto.Slug,
			Status:                       dto.Status,
			LastIncident:                 dto.LastIncident,
			LastIncidentLink:             dto.LastIncidentLink,
			IncludeAllComponents:         dto.IncludeAllComponents,
			MonitorIncidents:             dto.MonitorIncidents,
			MonitorScheduledMaintenances: dto.MonitorScheduledMaintenances,
			UpcomingMaintenance:          dto.UpcomingMaintenance,
			UpcomingMaintenanceLink:      dto.UpcomingMaintenanceLink,
			ComponentIDs:                 dto.ComponentIDs,
			ComponentGroupIDs:            dto.ComponentGroupIDs,
		}
	}

	return result, totalCount, upCount, downCount, nil
}
