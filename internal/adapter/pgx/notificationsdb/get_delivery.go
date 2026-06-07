package notificationsdb

import (
	"context"
	_ "embed"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/internal/common/apperrors"
	"github.com/yash492/statusy/internal/domain/notifications"
)

//go:embed queries/get_notification_delivery.sql
var getNotificationDeliveryQuery string

// GetDelivery gets existing notification deliveries for target and alert
func (r *PostgresNotificationsRepository) GetDelivery(ctx context.Context, channelID uint, alertType notifications.AlertType, alertID uint) (notifications.NotificationDelivery, error) {
	var d notifications.NotificationDelivery
	err := r.readDB.QueryRow(ctx, getNotificationDeliveryQuery, channelID, alertType, alertID).Scan(
		&d.ID, &d.ViewNotificationID, &d.AlertType, &d.AlertID, &d.LastUpdateID, &d.ExternalIdentifier, &d.CreatedAt, &d.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return notifications.NotificationDelivery{}, apperrors.NotFoundError("notification delivery not found", err)
		}
		return notifications.NotificationDelivery{}, apperrors.InternalError("failed to fetch notification delivery", err)
	}
	return d, nil
}
