package notificationsdb

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/internal/domain/notifications"
)

//go:embed queries/get_maintenance_notification_details.sql
var getMaintenanceNotificationDetailsQuery string

type maintenanceDetailsDto struct {
	MaintenanceID uint      `db:"maintenance_id"`
	UpdateID      uint      `db:"update_id"`
	Title         string    `db:"title"`
	Status        string    `db:"status"`
	Description   string    `db:"description"`
	ProviderID    string    `db:"provider_id"`
	ServiceName   string    `db:"service_name"`
	Components    []byte    `db:"components"`
	StartTime     time.Time `db:"start_time"`
	EndTime       time.Time `db:"end_time"`
	UpdatedAt     time.Time `db:"updated_at"`
	Link          string    `db:"link"`
}

// GetMaintenanceNotificationDetails gets all notification details for maintenance update
func (r *PostgresNotificationsRepository) GetMaintenanceNotificationDetails(ctx context.Context, updateID uint) (notifications.MaintenanceNotificationDetails, error) {
	rows, err := r.readDB.Query(ctx, getMaintenanceNotificationDetailsQuery, updateID)
	if err != nil {
		return notifications.MaintenanceNotificationDetails{}, fmt.Errorf("failed to get maintenance notification details: %w", err)
	}
	defer rows.Close()

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[maintenanceDetailsDto])
	if err != nil {
		return notifications.MaintenanceNotificationDetails{}, fmt.Errorf("failed to collect maintenance notification details: %w", err)
	}

	var components []notifications.NotificationComponent
	if err := json.Unmarshal(item.Components, &components); err != nil {
		return notifications.MaintenanceNotificationDetails{}, fmt.Errorf("failed to unmarshal maintenance components: %w", err)
	}

	return notifications.MaintenanceNotificationDetails{
		MaintenanceID: item.MaintenanceID,
		UpdateID:      item.UpdateID,
		Title:         item.Title,
		Status:        item.Status,
		Description:   item.Description,
		ProviderID:    item.ProviderID,
		ServiceName:   item.ServiceName,
		Components:    components,
		StartTime:     item.StartTime,
		EndTime:       item.EndTime,
		UpdatedAt:     item.UpdatedAt,
		Link:          item.Link,
	}, nil
}
