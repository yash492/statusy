package notificationsdb

import (
	"context"
	_ "embed"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/internal/common/apperrors"
	"github.com/yash492/statusy/internal/domain/notifications"
)

//go:embed queries/get_notification_delivery.sql
var getNotificationDeliveryQuery string

type notificationDeliveryDto struct {
	ID                 uint                    `db:"id"`
	ViewNotificationID uint                    `db:"view_notification_id"`
	AlertType          notifications.AlertType `db:"alert_type"`
	AlertID            uint                    `db:"alert_id"`
	LastUpdateID       uint                    `db:"last_update_id"`
	ExternalIdentifier string                  `db:"external_identifier"`
	CreatedAt          time.Time               `db:"created_at"`
	UpdatedAt          time.Time               `db:"updated_at"`
}

// GetDelivery gets existing notification deliveries for target and alert
func (r *PostgresNotificationsRepository) GetDelivery(ctx context.Context, channelID uint, alertType notifications.AlertType, alertID uint) (notifications.NotificationDelivery, error) {
	rows, err := r.readDB.Query(ctx, getNotificationDeliveryQuery, channelID, alertType, alertID)
	if err != nil {
		return notifications.NotificationDelivery{}, apperrors.InternalError("failed to fetch notification delivery", err)
	}
	defer rows.Close()

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[notificationDeliveryDto])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return notifications.NotificationDelivery{}, apperrors.NotFoundError("notification delivery not found", err)
		}
		return notifications.NotificationDelivery{}, apperrors.InternalError("failed to collect notification delivery", err)
	}

	return notifications.NotificationDelivery{
		ID:                 item.ID,
		ViewNotificationID: item.ViewNotificationID,
		AlertType:          item.AlertType,
		AlertID:            item.AlertID,
		LastUpdateID:       item.LastUpdateID,
		ExternalIdentifier: item.ExternalIdentifier,
		CreatedAt:          item.CreatedAt,
		UpdatedAt:          item.UpdatedAt,
	}, nil
}
