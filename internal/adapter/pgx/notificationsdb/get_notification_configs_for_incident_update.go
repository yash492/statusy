package notificationsdb

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/yash492/statusy/internal/domain/notifications"
)

//go:embed queries/get_notification_configs_for_incident_update.sql
var getNotificationConfigsForIncidentUpdateQuery string

// GetNotificationConfigsForIncidentUpdate retrieves all view notification configurations that should receive updates for this incident update
func (r *PostgresNotificationsRepository) GetNotificationConfigsForIncidentUpdate(ctx context.Context, updateID uint) ([]notifications.ViewNotification, error) {
	rows, err := r.readDB.Query(ctx, getNotificationConfigsForIncidentUpdateQuery, updateID)
	if err != nil {
		return nil, fmt.Errorf("failed to get notification configs for incident update: %w", err)
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
