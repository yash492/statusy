package notificationsdb

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/yash492/statusy/internal/domain/notifications"
)

//go:embed queries/get_view_notifications_by_view_id.sql
var getViewNotificationsByViewIDQuery string

// GetByViewID returns all notification destinations config for a given view
func (r *PostgresNotificationsRepository) GetByViewID(ctx context.Context, viewID uint) ([]notifications.ViewNotification, error) {
	rows, err := r.readDB.Query(ctx, getViewNotificationsByViewIDQuery, viewID)
	if err != nil {
		return nil, fmt.Errorf("failed to query view notifications: %w", err)
	}
	defer rows.Close()

	var list []notifications.ViewNotification
	for rows.Next() {
		var vn notifications.ViewNotification
		err := rows.Scan(&vn.ID, &vn.ViewID, &vn.Type, &vn.Config, &vn.CreatedAt, &vn.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan view notification: %w", err)
		}
		list = append(list, vn)
	}
	return list, nil
}
