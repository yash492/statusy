package notificationsdb

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/yash492/statusy/internal/domain/notifications"
)

//go:embed queries/save_view_notification.sql
var saveViewNotificationQuery string

// Save view notifications
func (r *PostgresNotificationsRepository) Save(ctx context.Context, vn notifications.ViewNotification) (notifications.ViewNotification, error) {
	var res notifications.ViewNotification
	err := r.writeDB.QueryRow(ctx, saveViewNotificationQuery, vn.ViewID, vn.Type, vn.Config).Scan(
		&res.ID, &res.ViewID, &res.Type, &res.Config, &res.CreatedAt, &res.UpdatedAt,
	)
	if err != nil {
		return notifications.ViewNotification{}, fmt.Errorf("failed to save view notification: %w", err)
	}
	return res, nil
}
