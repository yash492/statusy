package notificationsdb

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/yash492/statusy/internal/domain/notifications"
)

//go:embed queries/get_maintenance_notification_details.sql
var getMaintenanceNotificationDetailsQuery string

// GetMaintenanceNotificationDetails gets all notification details for maintenance update
func (r *PostgresNotificationsRepository) GetMaintenanceNotificationDetails(ctx context.Context, updateID uint) (notifications.MaintenanceNotificationDetails, error) {
	var d notifications.MaintenanceNotificationDetails
	err := r.readDB.QueryRow(ctx, getMaintenanceNotificationDetailsQuery, updateID).Scan(
		&d.MaintenanceID, &d.UpdateID, &d.Title, &d.Status, &d.Description, &d.ProviderID, &d.ServiceName, &d.ComponentNames, &d.StartTime, &d.EndTime, &d.UpdatedAt,
	)
	if err != nil {
		return notifications.MaintenanceNotificationDetails{}, fmt.Errorf("failed to get maintenance notification details: %w", err)
	}
	return d, nil
}
