package notification

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/yash492/statusy/internal/domain/notifications"
	"resty.dev/v3"
)

type PagerDutyConfig struct {
	RoutingKey string `json:"routing_key"`
}

type PagerDutyDispatcher struct {
	client *resty.Client
}

func newPagerDutyDispatcher(client *resty.Client) *PagerDutyDispatcher {
	return &PagerDutyDispatcher{client: client}
}

// Ensure PagerDutyDispatcher implements ChannelDispatcher interface
var _ ChannelDispatcher = &PagerDutyDispatcher{}

func (p *PagerDutyDispatcher) Send(
	ctx context.Context,
	configRaw json.RawMessage,
	isFirst bool,
	isResolve bool,
	data AlertData,
	prevExtID string,
) (string, error) {
	var cfg PagerDutyConfig
	if err := json.Unmarshal(configRaw, &cfg); err != nil {
		return "", fmt.Errorf("failed to parse PagerDuty config: %w", err)
	}
	if cfg.RoutingKey == "" {
		return "", fmt.Errorf("PagerDuty routing key is empty")
	}

	eventAction := "trigger"
	severity := "critical"
	if isResolve {
		eventAction = "resolve"
	}
	if data.AlertType == notifications.AlertTypeScheduledMaintenance {
		severity = "info"
	}

	dedupKey := fmt.Sprintf("statusy-%s-%d", string(data.AlertType), data.AlertID)
	comps := formatComponents(data.Components)

	payload := pagerduty.V2Event{
		RoutingKey: cfg.RoutingKey,
		Action:     eventAction,
		DedupKey:   dedupKey,
		Payload: &pagerduty.V2Payload{
			Summary:  fmt.Sprintf("[%s] %s", data.ServiceName, data.Title),
			Source:   "Statusy",
			Severity: severity,
			Details: map[string]any{
				"status":              data.Status,
				"description":         data.Description,
				"affected_components": comps,
				"link":                data.Link,
			},
		},
	}

	resp, err := p.client.R().
		SetContext(ctx).
		SetBody(payload).
		Post("https://events.pagerduty.com/v2/enqueue")
	if err != nil {
		return "", fmt.Errorf("failed to send PagerDuty HTTP request: %w", err)
	}
	if resp.IsError() {
		return "", fmt.Errorf("PagerDuty returned status code %d: %s", resp.StatusCode(), resp.String())
	}

	return dedupKey, nil
}
