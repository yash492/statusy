package notificationsdb

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/yash492/statusy/internal/domain/notifications"
)

//go:embed queries/update_view_notification.sql
var updateViewNotificationQuery string

// Update updates an existing view notification config
func (r *PostgresNotificationsRepository) Update(ctx context.Context, vn notifications.ViewNotification) (notifications.ViewNotification, error) {
	var res notifications.ViewNotification
	err := r.writeDB.QueryRow(ctx, updateViewNotificationQuery, vn.ID, vn.Name, vn.Type, vn.Config).Scan(
		&res.ID, &res.ViewID, &res.Name, &res.Type, &res.Config, &res.CreatedAt, &res.UpdatedAt,
	)
	if err != nil {
		return notifications.ViewNotification{}, fmt.Errorf("failed to update view notification: %w", err)
	}
	return res, nil
}
