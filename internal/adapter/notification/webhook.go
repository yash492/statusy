package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/yash492/statusy/internal/domain/notifications"
	"resty.dev/v3"
)

type WebhookConfig struct {
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}

type customWebhookPayload struct {
	EventType   string                                `json:"event_type"`
	AlertID     uint                                  `json:"alert_id"`
	UpdateID    uint                                  `json:"update_id"`
	Title       string                                `json:"title"`
	Status      string                                `json:"status"`
	Description string                                `json:"description"`
	ServiceName string                                `json:"service_name"`
	Components  []notifications.NotificationComponent `json:"components"`
	Link        string                                `json:"link"`
	UpdatedAt   string                                `json:"updated_at"`
	StartTime   *string                               `json:"start_time,omitempty"`
	EndTime     *string                               `json:"end_time,omitempty"`
}

type WebhookDispatcher struct {
	client *resty.Client
}

func newWebhookDispatcher(client *resty.Client) *WebhookDispatcher {
	return &WebhookDispatcher{client: client}
}

// Ensure WebhookDispatcher implements ChannelDispatcher interface
var _ ChannelDispatcher = &WebhookDispatcher{}

func (w *WebhookDispatcher) Send(
	ctx context.Context,
	configRaw json.RawMessage,
	isFirst bool,
	isResolve bool,
	data AlertData,
	prevExtID string,
) (string, error) {
	var cfg WebhookConfig
	if err := json.Unmarshal(configRaw, &cfg); err != nil {
		return "", fmt.Errorf("failed to parse custom Webhook config: %w", err)
	}
	if cfg.URL == "" {
		return "", fmt.Errorf("custom Webhook URL is empty")
	}

	var startTimeStr, endTimeStr *string
	if data.StartTime != nil {
		s := data.StartTime.UTC().Format(time.RFC3339)
		startTimeStr = &s
	}
	if data.EndTime != nil {
		e := data.EndTime.UTC().Format(time.RFC3339)
		endTimeStr = &e
	}

	payload := customWebhookPayload{
		EventType:   data.AlertType,
		AlertID:     data.AlertID,
		UpdateID:    data.UpdateID,
		Title:       data.Title,
		Status:      data.Status,
		Description: data.Description,
		ServiceName: data.ServiceName,
		Components:  data.Components,
		Link:        data.Link,
		UpdatedAt:   data.UpdatedAt.UTC().Format(time.RFC3339),
		StartTime:   startTimeStr,
		EndTime:     endTimeStr,
	}

	req := w.client.R().
		SetContext(ctx).
		SetBody(payload)

	for k, v := range cfg.Headers {
		req.SetHeader(k, v)
	}

	resp, err := req.Post(cfg.URL)
	if err != nil {
		return "", fmt.Errorf("failed to send custom Webhook HTTP request: %w", err)
	}
	if resp.IsError() {
		return "", fmt.Errorf("custom Webhook returned status code %d: %s", resp.StatusCode(), resp.String())
	}

	return "", nil
}
