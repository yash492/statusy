package notificationsdb

import (
	"context"
	_ "embed"

	"github.com/yash492/statusy/internal/domain/notifications"
)

//go:embed queries/get_notification_delivery.sql
var getNotificationDeliveryQuery string

// GetDelivery gets existing notification deliveries for target and alert
func (r *PostgresNotificationsRepository) GetDelivery(ctx context.Context, channelID uint, alertType string, alertID uint) (notifications.NotificationDelivery, error) {
	var d notifications.NotificationDelivery
	err := r.readDB.QueryRow(ctx, getNotificationDeliveryQuery, channelID, alertType, alertID).Scan(
		&d.ID, &d.ViewNotificationID, &d.AlertType, &d.AlertID, &d.LastUpdateID, &d.ExternalIdentifier, &d.CreatedAt, &d.UpdatedAt,
	)
	if err != nil {
		return notifications.NotificationDelivery{}, err
	}
	return d, nil
}
