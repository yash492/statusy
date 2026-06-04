package incidentsdb

import (
	"context"
	_ "embed"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/internal/common/apperrors"
)

//go:embed queries/update_resolution_status.sql
var updateResolutionStatusQuery string

func (c *PostgresIncidentRepository) UpdateResolutionStatus(ctx context.Context, serviceIDs []uint) error {
	_, err := c.writeDB.Exec(ctx, updateResolutionStatusQuery, pgx.NamedArgs{
		"service_ids": serviceIDs,
	})
	if err != nil {
		c.lg.ErrorContext(ctx, "error updating incident resolution status", slog.Any("service_ids", serviceIDs), slog.Any("err", err))
		return apperrors.InternalError("failed to update incident resolution status", err)
	}

	c.lg.InfoContext(ctx, "incident resolution status updated", slog.Int("service_count", len(serviceIDs)))
	return nil
}
