package notificationsdb

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/yash492/statusy/internal/domain/notifications"
)

//go:embed queries/save_delivery_failure.sql
var saveDeliveryFailureQuery string

// SaveDeliveryFailure logs notification delivery failure in database
func (r *PostgresNotificationsRepository) SaveDeliveryFailure(ctx context.Context, failure notifications.NotificationDeliveryFailure) error {
	_, err := r.writeDB.Exec(ctx, saveDeliveryFailureQuery, failure.ViewNotificationID, failure.AlertType, failure.AlertID, failure.UpdateID, failure.ErrorMessage)
	if err != nil {
		return fmt.Errorf("failed to save notification delivery failure: %w", err)
	}
	return nil
}
