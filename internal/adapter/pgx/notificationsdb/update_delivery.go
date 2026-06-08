package notificationsdb

import (
	"context"
	_ "embed"
	"fmt"
)

//go:embed queries/update_notification_delivery.sql
var updateNotificationDeliveryQuery string

// UpdateDelivery updates existing delivery
func (r *PostgresNotificationsRepository) UpdateDelivery(ctx context.Context, deliveryID uint, lastUpdateID uint, externalIdentifier string) error {
	_, err := r.writeDB.Exec(ctx, updateNotificationDeliveryQuery, lastUpdateID, externalIdentifier, deliveryID)
	if err != nil {
		return fmt.Errorf("failed to update notification delivery: %w", err)
	}
	return nil
}
