package notificationsdb

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/yash492/statusy/internal/domain/notifications"
)

//go:embed queries/save_notification_delivery.sql
var saveNotificationDeliveryQuery string

// SaveDelivery saves delivery mapping
func (r *PostgresNotificationsRepository) SaveDelivery(ctx context.Context, delivery notifications.NotificationDelivery) error {
	_, err := r.writeDB.Exec(ctx, saveNotificationDeliveryQuery, delivery.ViewNotificationID, delivery.AlertType, delivery.AlertID, delivery.LastUpdateID, delivery.ExternalIdentifier)
	if err != nil {
		return fmt.Errorf("failed to save notification delivery: %w", err)
	}
	return nil
}
