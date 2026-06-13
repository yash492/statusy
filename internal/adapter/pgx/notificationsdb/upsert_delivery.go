package notificationsdb

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/yash492/statusy/internal/domain/notifications"
)

//go:embed queries/save_notification_delivery.sql
var upsertNotificationDeliveryQuery string

func (r *PostgresNotificationsRepository) UpsertDelivery(ctx context.Context, delivery notifications.NotificationDelivery) error {
	_, err := r.writeDB.Exec(ctx, upsertNotificationDeliveryQuery, delivery.ViewNotificationID, delivery.AlertType, delivery.AlertID, delivery.LastUpdateID, delivery.ExternalIdentifier)
	if err != nil {
		return fmt.Errorf("failed to upsert notification delivery: %w", err)
	}
	return nil
}
