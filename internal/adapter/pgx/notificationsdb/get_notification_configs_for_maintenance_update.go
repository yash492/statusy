package notificationsdb

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/yash492/statusy/internal/domain/notifications"
)

//go:embed queries/get_notification_configs_for_maintenance_update.sql
var getNotificationConfigsForMaintenanceUpdateQuery string

// GetNotificationConfigsForMaintenanceUpdate retrieves all view notification configurations that should receive updates for this scheduled maintenance update
func (r *PostgresNotificationsRepository) GetNotificationConfigsForMaintenanceUpdate(ctx context.Context, updateID uint) ([]notifications.ViewNotification, error) {
	rows, err := r.readDB.Query(ctx, getNotificationConfigsForMaintenanceUpdateQuery, updateID)
	if err != nil {
		return nil, fmt.Errorf("failed to get notification configs for maintenance update: %w", err)
	}
	defer rows.Close()

	var list []notifications.ViewNotification
	for rows.Next() {
		var vn notifications.ViewNotification
		err := rows.Scan(&vn.ID, &vn.ViewID, &vn.Type, &vn.Config, &vn.CreatedAt, &vn.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan resolved view notification config: %w", err)
		}
		list = append(list, vn)
	}
	return list, nil
}
