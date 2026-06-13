package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/yash492/statusy/internal/common/jsonutil"
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
	data AlertData,
	prevExtID string,
) (string, error) {
	cfg, err := jsonutil.UnmarshalWithType[PagerDutyConfig](configRaw)
	if err != nil {
		return "", err
	}
	if cfg.RoutingKey == "" {
		return "", fmt.Errorf("PagerDuty routing key is empty")
	}

	eventAction := "trigger"
	if data.Status.IsResolved() {
		eventAction = "resolve"
	}

	severity := "critical"
	if data.AlertType == notifications.AlertTypeScheduledMaintenance {
		severity = "info"
	}

	dedupKey := fmt.Sprintf("statusy-%s-%d", string(data.AlertType), data.AlertID)
	comps := formatComponents(data.Components)

	_, err = pagerduty.ManageEventWithContext(ctx, pagerduty.V2Event{
		RoutingKey: cfg.RoutingKey,
		Action:     eventAction,
		Client:     data.ServiceName,
		DedupKey:   dedupKey,
		ClientURL:  data.Link,
		Payload: &pagerduty.V2Payload{
			Timestamp: data.UpdatedAt.UTC().Format(time.RFC3339),
			Summary:   fmt.Sprintf("%s: %s", data.ServiceName, data.Title),
			Source:    "Statusy",
			Severity:  severity,
			Component: comps,
			Details: map[string]any{
				"status":      data.Status.String(),
				"description": data.Description,
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to send PagerDuty event: %w", err)
	}

	return dedupKey, nil
}
