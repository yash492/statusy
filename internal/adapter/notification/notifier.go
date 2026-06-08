package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/yash492/statusy/internal/domain/notifications"
	"resty.dev/v3"
)

type ChannelDispatcher interface {
	Send(ctx context.Context, config json.RawMessage, isFirst bool, isResolve bool, data AlertData, prevExtID string) (string, error)
}

type HttpNotifier struct {
	client      *resty.Client
	lg          *slog.Logger
	dispatchers map[notifications.NotificationType]ChannelDispatcher
}

func NewHttpNotifier(lg *slog.Logger) *HttpNotifier {
	client := resty.New()
	client.SetTimeout(10 * time.Second)
	client.SetRetryCount(3)
	client.SetRetryWaitTime(1 * time.Second)
	client.SetRetryMaxWaitTime(5 * time.Second)
	client.AddRetryConditions(func(resp *resty.Response, err error) bool {
		return err != nil || (resp != nil && resp.StatusCode() >= 500)
	})

	h := &HttpNotifier{
		client: client,
		lg:     lg,
	}

	h.dispatchers = map[notifications.NotificationType]ChannelDispatcher{
		notifications.NotificationTypeSlack:                      newSlackDispatcher(client),
		notifications.NotificationTypeDiscord:                    newDiscordDispatcher(client, lg),
		notifications.NotificationTypeMsTeams:                    newMsTeamsDispatcher(lg),
		notifications.NotificationTypePagerDuty:                  newPagerDutyDispatcher(client),
		notifications.NotificationTypeSolarwindsIncidentResponse: newSolarwindsDispatcher(client),
		notifications.NotificationTypeWebhook:                    newWebhookDispatcher(client),
	}

	return h
}

// Ensure HttpNotifier implements notifications.Notifier interface
var _ notifications.Notifier = &HttpNotifier{}

type AlertData struct {
	AlertID     uint
	UpdateID    uint
	Title       string
	Status      string
	Description string
	ServiceName string
	Components  []notifications.NotificationComponent
	Link        string
	UpdatedAt   time.Time
	StartTime   *time.Time
	EndTime     *time.Time
	AlertType   notifications.AlertType
}

func (h *HttpNotifier) SendIncident(
	ctx context.Context,
	ch notifications.ViewNotification,
	isFirst bool,
	details notifications.IncidentNotificationDetails,
	prevExtID string,
) (string, error) {
	data := AlertData{
		AlertID:     details.IncidentID,
		UpdateID:    details.UpdateID,
		Title:       details.Title,
		Status:      details.Status,
		Description: details.Description,
		ServiceName: details.ServiceName,
		Components:  details.Components,
		Link:        details.Link,
		UpdatedAt:   details.UpdatedAt,
		AlertType:   notifications.AlertTypeIncident,
	}
	return h.send(ctx, ch, isFirst, data, prevExtID)
}

func (h *HttpNotifier) SendMaintenance(
	ctx context.Context,
	ch notifications.ViewNotification,
	isFirst bool,
	details notifications.MaintenanceNotificationDetails,
	prevExtID string,
) (string, error) {
	data := AlertData{
		AlertID:     details.MaintenanceID,
		UpdateID:    details.UpdateID,
		Title:       details.Title,
		Status:      details.Status,
		Description: details.Description,
		ServiceName: details.ServiceName,
		Components:  details.Components,
		Link:        details.Link,
		UpdatedAt:   details.UpdatedAt,
		StartTime:   &details.StartTime,
		EndTime:     &details.EndTime,
		AlertType:   notifications.AlertTypeScheduledMaintenance,
	}
	return h.send(ctx, ch, isFirst, data, prevExtID)
}

func (h *HttpNotifier) send(
	ctx context.Context,
	ch notifications.ViewNotification,
	isFirst bool,
	data AlertData,
	prevExtID string,
) (string, error) {
	isResolve := isResolve(data.AlertType, data.Status)

	dispatcher, ok := h.dispatchers[ch.Type]
	if !ok {
		return "", fmt.Errorf("unsupported notification channel type: %s", ch.Type)
	}

	return dispatcher.Send(ctx, ch.Config, isFirst, isResolve, data, prevExtID)
}

func isResolve(alertType notifications.AlertType, status string) bool {
	if alertType == notifications.AlertTypeIncident {
		return strings.ToLower(status) == "resolved"
	}
	return strings.ToLower(status) == "completed"
}

func getColor(isFirst, isResolve bool) (string, int) {
	if isResolve {
		return "#2EB886", 3061894
	}
	if isFirst {
		return "#E03E3E", 14696000
	}
	return "#FFA500", 16753920
}

func formatComponents(components []notifications.NotificationComponent) string {
	if len(components) == 0 {
		return "None"
	}
	var compStrings []string
	for _, c := range components {
		if c.GroupName != nil && *c.GroupName != "" {
			compStrings = append(compStrings, fmt.Sprintf("%s: %s", *c.GroupName, c.Name))
		} else {
			compStrings = append(compStrings, c.Name)
		}
	}
	if len(compStrings) <= 5 {
		return strings.Join(compStrings, ", ")
	}
	return fmt.Sprintf("%s, and %d more", strings.Join(compStrings[:5], ", "), len(compStrings)-5)
}
