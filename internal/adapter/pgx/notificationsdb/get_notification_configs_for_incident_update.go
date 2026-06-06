package notificationsdb

import (
	"context"
	_ "embed"

	"github.com/yash492/statusy/internal/common/apperrors"
	"github.com/yash492/statusy/internal/domain/notifications"
)

//go:embed queries/get_notification_configs_for_incident_update.sql
var getNotificationConfigsForIncidentUpdateQuery string

// GetNotificationConfigsForIncidentUpdate retrieves all view notification configurations that should receive updates for this incident update
func (r *PostgresNotificationsRepository) GetNotificationConfigsForIncidentUpdate(ctx context.Context, updateID uint) ([]notifications.ViewNotification, error) {
	rows, err := r.readDB.Query(ctx, getNotificationConfigsForIncidentUpdateQuery, updateID)
	if err != nil {
		return nil, apperrors.InternalError("failed to get notification configs for incident update", err)
	}
	defer rows.Close()

	var list []notifications.ViewNotification
	for rows.Next() {
		var vn notifications.ViewNotification
		err := rows.Scan(&vn.ID, &vn.ViewID, &vn.Type, &vn.Config, &vn.CreatedAt, &vn.UpdatedAt)
		if err != nil {
			return nil, apperrors.InternalError("failed to scan resolved view notification config", err)
		}
		list = append(list, vn)
	}
	return list, nil
}
