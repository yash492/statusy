package notificationsdb

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/internal/domain/notifications"
)

//go:embed queries/update_view_notification.sql
var updateViewNotificationQuery string

// Update updates an existing view notification config
func (r *PostgresNotificationsRepository) Update(ctx context.Context, vn notifications.ViewNotification) (notifications.ViewNotification, error) {
	rows, err := r.writeDB.Query(ctx, updateViewNotificationQuery, vn.ID, vn.Name, vn.Type, vn.Config)
	if err != nil {
		return notifications.ViewNotification{}, fmt.Errorf("failed to update view notification: %w", err)
	}
	defer rows.Close()

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[viewNotificationDto])
	if err != nil {
		return notifications.ViewNotification{}, fmt.Errorf("failed to collect updated view notification: %w", err)
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
