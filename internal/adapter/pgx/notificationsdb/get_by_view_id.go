package notificationsdb

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/internal/domain/notifications"
)

//go:embed queries/get_view_notifications_by_view_id.sql
var getViewNotificationsByViewIDQuery string

//go:embed queries/count_view_notifications.sql
var countViewNotificationsQuery string

// GetByViewID returns all notification destinations config for a given view with search and pagination
func (r *PostgresNotificationsRepository) GetByViewID(ctx context.Context, viewID uint, search string, limit int, offset int) ([]notifications.ViewNotification, int64, error) {
	rows, err := r.readDB.Query(ctx, getViewNotificationsByViewIDQuery, pgx.NamedArgs{
		"view_id": viewID,
		"search":  search,
		"limit":   limit,
		"offset":  offset,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query view notifications: %w", err)
	}
	defer rows.Close()

	var list []notifications.ViewNotification
	for rows.Next() {
		var vn notifications.ViewNotification
		err := rows.Scan(&vn.ID, &vn.ViewID, &vn.Name, &vn.Type, &vn.Config, &vn.CreatedAt, &vn.UpdatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan view notification: %w", err)
		}
		list = append(list, vn)
	}

	var totalCount int64
	err = r.readDB.QueryRow(ctx, countViewNotificationsQuery, pgx.NamedArgs{
		"view_id": viewID,
		"search":  search,
	}).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count view notifications: %w", err)
	}

	return list, totalCount, nil
}
