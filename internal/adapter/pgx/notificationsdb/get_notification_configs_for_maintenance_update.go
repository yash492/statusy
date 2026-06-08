package notificationsdb

import (
	"context"
	_ "embed"

	"github.com/jackc/pgx/v5"
	"github.com/samber/lo"
	"github.com/yash492/statusy/internal/common/apperrors"
	"github.com/yash492/statusy/internal/domain/notifications"
)

//go:embed queries/get_notification_configs_for_maintenance_update.sql
var getNotificationConfigsForMaintenanceUpdateQuery string

// GetNotificationConfigsForMaintenanceUpdate retrieves all view notification configurations that should receive updates for this scheduled maintenance update
func (r *PostgresNotificationsRepository) GetNotificationConfigsForMaintenanceUpdate(ctx context.Context, updateID uint) ([]notifications.ViewNotification, error) {
	rows, err := r.readDB.Query(ctx, getNotificationConfigsForMaintenanceUpdateQuery, updateID)
	if err != nil {
		return nil, apperrors.InternalError("failed to get notification configs for maintenance update", err)
	}
	defer rows.Close()

	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[viewNotificationDto])
	if err != nil {
		return nil, apperrors.InternalError("failed to collect notification configs for maintenance update", err)
	}

	return lo.Map(dtos, func(item viewNotificationDto, _ int) notifications.ViewNotification {
		return notifications.ViewNotification{
			ID:        item.ID,
			ViewID:    item.ViewID,
			Name:      item.Name,
			Type:      item.Type,
			Config:    item.Config,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		}
	}), nil
}
