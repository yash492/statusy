package notificationsdb

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/internal/domain/notifications"
)

//go:embed queries/save_view_notification.sql
var saveViewNotificationQuery string

// Save view notifications
func (r *PostgresNotificationsRepository) Save(ctx context.Context, vn notifications.ViewNotification) (notifications.ViewNotification, error) {
	rows, err := r.writeDB.Query(ctx, saveViewNotificationQuery, vn.ViewID, vn.Name, vn.Type, vn.Config)
	if err != nil {
		return notifications.ViewNotification{}, fmt.Errorf("failed to save view notification: %w", err)
	}
	defer rows.Close()

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[viewNotificationDto])
	if err != nil {
		return notifications.ViewNotification{}, fmt.Errorf("failed to collect saved view notification: %w", err)
	}

	return notifications.ViewNotification{
		ID:        item.ID,
		ViewID:    item.ViewID,
		Name:      item.Name,
		Type:      item.Type,
		Config:    item.Config,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}, nil
}
