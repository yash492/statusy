package notificationsdb

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/samber/lo"
	"github.com/yash492/statusy/internal/domain/notifications"
)

//go:embed queries/get_view_notifications_by_view_id.sql
var getViewNotificationsByViewIDQuery string

//go:embed queries/count_view_notifications.sql
var countViewNotificationsQuery string

type viewNotificationDto struct {
	ID        uint                           `db:"id"`
	ViewID    uint                           `db:"view_id"`
	Name      string                         `db:"name"`
	Type      notifications.NotificationType `db:"type"`
	Config    json.RawMessage                `db:"config"`
	CreatedAt time.Time                      `db:"created_at"`
	UpdatedAt time.Time                      `db:"updated_at"`
}

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

	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[viewNotificationDto])
	if err != nil {
		return nil, 0, fmt.Errorf("failed to collect view notification rows: %w", err)
	}

	var totalCount int64
	err = r.readDB.QueryRow(ctx, countViewNotificationsQuery, pgx.NamedArgs{
		"view_id": viewID,
		"search":  search,
	}).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count view notifications: %w", err)
	}

	result := lo.Map(dtos, func(item viewNotificationDto, _ int) notifications.ViewNotification {
		return notifications.ViewNotification{
			ID:        item.ID,
			ViewID:    item.ViewID,
			Name:      item.Name,
			Type:      item.Type,
			Config:    item.Config,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		}
	})

	return result, totalCount, nil
}
