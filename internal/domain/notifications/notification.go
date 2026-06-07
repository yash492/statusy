package notifications

import (
	"context"
	"encoding/json"
	"time"
)

type NotificationType string

const (
	NotificationTypeSlack                      NotificationType = "slack"
	NotificationTypeDiscord                    NotificationType = "discord"
	NotificationTypeMsTeams                    NotificationType = "msteams"
	NotificationTypePagerDuty                  NotificationType = "pagerduty"
	NotificationTypeSolarwindsIncidentResponse NotificationType = "solarwinds_incident_response"
	NotificationTypeWebhook                    NotificationType = "webhook"
)

type ViewNotification struct {
	ID        uint             `json:"id"`
	ViewID    uint             `json:"view_id"`
	Name      string           `json:"name"`
	Type      NotificationType `json:"type"`
	Config    json.RawMessage  `json:"config"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

type NotificationDelivery struct {
	ID                 uint      `json:"id"`
	ViewNotificationID uint      `json:"view_notification_id"`
	AlertType          string    `json:"alert_type"` // "incident" or "scheduled_maintenance"
	AlertID            uint      `json:"alert_id"`
	LastUpdateID       uint      `json:"last_update_id"`
	ExternalIdentifier string    `json:"external_identifier"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type NotificationComponent struct {
	Name      string  `json:"name"`
	GroupName *string `json:"group_name"`
}

type NotificationDeliveryFailure struct {
	ID                 uint      `json:"id"`
	ViewNotificationID uint      `json:"view_notification_id"`
	AlertType          string    `json:"alert_type"` // "incident" or "scheduled_maintenance"
	AlertID            uint      `json:"alert_id"`
	UpdateID           uint      `json:"update_id"`
	ErrorMessage       string    `json:"error_message"`
	CreatedAt          time.Time `json:"created_at"`
}

// IncidentNotificationDetails holds detailed info of an incident update to build notification payload
type IncidentNotificationDetails struct {
	IncidentID  uint                    `json:"incident_id"`
	UpdateID    uint                    `json:"update_id"`
	Title       string                  `json:"title"`
	Status      string                  `json:"status"` // investigating, identified, monitoring, resolved
	Description string                  `json:"description"`
	ProviderID  string                  `json:"provider_id"`
	ServiceName string                  `json:"service_name"`
	Components  []NotificationComponent `json:"components"`
	UpdatedAt   time.Time               `json:"updated_at"`
	Link        string                  `json:"link"`
}

// MaintenanceNotificationDetails holds detailed info of a maintenance update to build notification payload
type MaintenanceNotificationDetails struct {
	MaintenanceID uint                    `json:"maintenance_id"`
	UpdateID      uint                    `json:"update_id"`
	Title         string                  `json:"title"`
	Status        string                  `json:"status"` // scheduled, in_progress, verifying, completed
	Description   string                  `json:"description"`
	ProviderID    string                  `json:"provider_id"`
	ServiceName   string                  `json:"service_name"`
	Components    []NotificationComponent `json:"components"`
	StartTime     time.Time               `json:"start_time"`
	EndTime       time.Time               `json:"end_time"`
	UpdatedAt     time.Time               `json:"updated_at"`
	Link          string                  `json:"link"`
}

type NotificationsRepository interface {
	Save(ctx context.Context, vn ViewNotification) (ViewNotification, error)
	GetByViewID(ctx context.Context, viewID uint, search string, limit int, offset int) ([]ViewNotification, int64, error)
	Delete(ctx context.Context, id uint) error
	Update(ctx context.Context, vn ViewNotification) (ViewNotification, error)

	GetDelivery(ctx context.Context, channelID uint, alertType string, alertID uint) (NotificationDelivery, error)
	SaveDelivery(ctx context.Context, delivery NotificationDelivery) error
	UpdateDelivery(ctx context.Context, deliveryID uint, lastUpdateID uint, externalIdentifier string) error
	SaveDeliveryFailure(ctx context.Context, failure NotificationDeliveryFailure) error

	GetNotificationConfigsForIncidentUpdate(ctx context.Context, updateID uint) ([]ViewNotification, error)
	GetNotificationConfigsForMaintenanceUpdate(ctx context.Context, updateID uint) ([]ViewNotification, error)
	GetIncidentNotificationDetails(ctx context.Context, updateID uint) (IncidentNotificationDetails, error)
	GetMaintenanceNotificationDetails(ctx context.Context, updateID uint) (MaintenanceNotificationDetails, error)
}

type Notifier interface {
	SendIncident(ctx context.Context, ch ViewNotification, isFirst bool, details IncidentNotificationDetails, prevExtID string) (string, error)
	SendMaintenance(ctx context.Context, ch ViewNotification, isFirst bool, details MaintenanceNotificationDetails, prevExtID string) (string, error)
}
