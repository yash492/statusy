package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/yash492/statusy/internal/common/jsonutil"
	"resty.dev/v3"
)

type SolarwindsConfig struct {
	WebhookURL string `json:"webhook_url"`
}

type squadcastPayload struct {
	Status      string `json:"status"`
	EventID     string `json:"event_id"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

type SolarwindsDispatcher struct {
	client *resty.Client
}

func newSolarwindsDispatcher(client *resty.Client) *SolarwindsDispatcher {
	return &SolarwindsDispatcher{client: client}
}

// Ensure SolarwindsDispatcher implements ChannelDispatcher interface
var _ ChannelDispatcher = &SolarwindsDispatcher{}

func (s *SolarwindsDispatcher) Send(
	ctx context.Context,
	configRaw json.RawMessage,
	isFirst bool,
	isResolve bool,
	data AlertData,
	prevExtID string,
) (string, error) {
	cfg, err := jsonutil.UnmarshalWithType[SolarwindsConfig](configRaw)
	if err != nil {
		return "", err
	}
	if cfg.WebhookURL == "" {
		return "", fmt.Errorf("SolarWinds Incident Response webhook URL is empty")
	}

	status := "trigger"
	if isResolve {
		status = "resolve"
	}

	payload := squadcastPayload{
		Status:      status,
		EventID:     strconv.FormatUint(uint64(data.AlertID), 10),
		Message:     fmt.Sprintf("[%s] Alert: %s", data.ServiceName, data.Title),
		Description: data.Description,
	}

	resp, err := s.client.R().
		SetContext(ctx).
		SetBody(payload).
		Post(cfg.WebhookURL)
	if err != nil {
		return "", fmt.Errorf("failed to send SolarWinds HTTP request: %w", err)
	}
	if resp.IsError() {
		return "", fmt.Errorf("SolarWinds returned status code %d: %s", resp.StatusCode(), resp.String())
	}

	return "", nil
}
