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

//go:embed queries/get_incident_notification_details.sql
var getIncidentNotificationDetailsQuery string

type incidentDetailsDto struct {
	IncidentID  uint      `db:"incident_id"`
	UpdateID    uint      `db:"update_id"`
	Title       string    `db:"title"`
	Status      string    `db:"status"`
	Description string    `db:"description"`
	ProviderID  string    `db:"provider_id"`
	ServiceName string    `db:"service_name"`
	Components  []byte    `db:"components"`
	UpdatedAt   time.Time `db:"updated_at"`
	Link        string    `db:"link"`
}

// GetIncidentNotificationDetails gets all notification details for incident update
func (r *PostgresNotificationsRepository) GetIncidentNotificationDetails(ctx context.Context, updateID uint) (notifications.IncidentNotificationDetails, error) {
	rows, err := r.readDB.Query(ctx, getIncidentNotificationDetailsQuery, updateID)
	if err != nil {
		return notifications.IncidentNotificationDetails{}, fmt.Errorf("failed to get incident notification details: %w", err)
	}
	defer rows.Close()

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[incidentDetailsDto])
	if err != nil {
		return notifications.IncidentNotificationDetails{}, fmt.Errorf("failed to collect incident notification details: %w", err)
	}

	var components []notifications.NotificationComponent
	if err := json.Unmarshal(item.Components, &components); err != nil {
		return notifications.IncidentNotificationDetails{}, fmt.Errorf("failed to unmarshal incident components: %w", err)
	}

	return notifications.IncidentNotificationDetails{
		IncidentID:  item.IncidentID,
		UpdateID:    item.UpdateID,
		Title:       item.Title,
		Status:      item.Status,
		Description: item.Description,
		ProviderID:  item.ProviderID,
		ServiceName: item.ServiceName,
		Components:  components,
		UpdatedAt:   item.UpdatedAt,
		Link:        item.Link,
	}, nil
}
