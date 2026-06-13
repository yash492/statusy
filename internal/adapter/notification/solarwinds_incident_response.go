package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/yash492/statusy/internal/common/jsonutil"
	"resty.dev/v3"
)

type SolarwindsConfig struct {
	WebhookURL string `json:"webhook_url"`
}

type squadcastIncidentEvent struct {
	Message     string       `json:"message"`
	Description string       `json:"description"`
	Tags        incidentTags `json:"tags"`
	Status      string       `json:"status"`
	EventID     string       `json:"event_id"`
}

type incidentTags map[string]incidentTag

type incidentTag struct {
	Value string `json:"value"`
	Color string `json:"color,omitempty"`
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
	if data.Status.IsResolved() {
		status = "resolve"
	}

	comps := formatComponents(data.Components)
	description := fmt.Sprintf("Updated At: %s\n", data.UpdatedAt.UTC().Format(time.RFC822)) +
		fmt.Sprintf("Affected Components: %s\n", comps) +
		fmt.Sprintf("Link: %s\n", data.Link) +
		data.Description

	payload := squadcastIncidentEvent{
		Message:     fmt.Sprintf("%s: %s", data.ServiceName, data.Title),
		Description: description,
		Tags: incidentTags{
			"Service":    {Value: data.ServiceName},
			"Components": {Value: comps},
			"Link":       {Value: data.Link},
		},
		Status:  status,
		EventID: strconv.FormatUint(uint64(data.AlertID), 10),
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
