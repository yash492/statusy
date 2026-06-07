package notificationsdb

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/yash492/statusy/internal/domain/notifications"
)

//go:embed queries/get_incident_notification_details.sql
var getIncidentNotificationDetailsQuery string

// GetIncidentNotificationDetails gets all notification details for incident update
func (r *PostgresNotificationsRepository) GetIncidentNotificationDetails(ctx context.Context, updateID uint) (notifications.IncidentNotificationDetails, error) {
	var d notifications.IncidentNotificationDetails
	var componentsJSON []byte
	err := r.readDB.QueryRow(ctx, getIncidentNotificationDetailsQuery, updateID).Scan(
		&d.IncidentID, &d.UpdateID, &d.Title, &d.Status, &d.Description, &d.ProviderID, &d.ServiceName, &componentsJSON, &d.UpdatedAt, &d.Link,
	)
	if err != nil {
		return notifications.IncidentNotificationDetails{}, fmt.Errorf("failed to get incident notification details: %w", err)
	}
	if err := json.Unmarshal(componentsJSON, &d.Components); err != nil {
		return notifications.IncidentNotificationDetails{}, fmt.Errorf("failed to unmarshal incident components details: %w", err)
	}
	return d, nil
}
