package notification

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/slack-go/slack"
	"resty.dev/v3"
)

type SlackConfig struct {
	WebhookURL string `json:"webhook_url"`
}

type SlackDispatcher struct {
	client *resty.Client
}

func newSlackDispatcher(client *resty.Client) *SlackDispatcher {
	return &SlackDispatcher{client: client}
}

// Ensure SlackDispatcher implements ChannelDispatcher interface
var _ ChannelDispatcher = &SlackDispatcher{}

func (s *SlackDispatcher) Send(
	ctx context.Context,
	configRaw json.RawMessage,
	isFirst bool,
	isResolve bool,
	data AlertData,
	prevExtID string,
) (string, error) {
	var cfg SlackConfig
	if err := json.Unmarshal(configRaw, &cfg); err != nil {
		return "", fmt.Errorf("failed to parse Slack config: %w", err)
	}
	if cfg.WebhookURL == "" {
		return "", fmt.Errorf("Slack webhook URL is empty")
	}

	color, _ := getColor(isFirst, isResolve)
	titleLink := fmt.Sprintf("*<%s|%s - %s>*", data.Link, data.ServiceName, data.Title)
	if data.Link == "" {
		titleLink = fmt.Sprintf("*%s - %s*", data.ServiceName, data.Title)
	}
	comps := formatComponents(data.Components)

	attachment := slack.Attachment{
		Color: color,
		Blocks: slack.Blocks{
			BlockSet: []slack.Block{
				slack.NewSectionBlock(
					slack.NewTextBlockObject("mrkdwn", titleLink, false, false),
					nil, nil,
				),
				slack.NewSectionBlock(
					slack.NewTextBlockObject("mrkdwn", data.Description, false, false),
					nil, nil,
				),
				slack.NewSectionBlock(
					nil,
					[]*slack.TextBlockObject{
						slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Status:*\n%s", data.Status), false, false),
						slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Affected Components:*\n%s", comps), false, false),
					},
					nil,
				),
				slack.NewContextBlock(
					"",
					slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("Statusy | %s", data.UpdatedAt.UTC().Format("2006-01-02 15:04:05 MST")), false, false),
				),
			},
		},
	}

	payload := slack.WebhookMessage{
		Attachments: []slack.Attachment{attachment},
	}

	resp, err := s.client.R().
		SetContext(ctx).
		SetBody(payload).
		Post(cfg.WebhookURL)
	if err != nil {
		return "", fmt.Errorf("failed to send Slack HTTP request: %w", err)
	}
	if resp.IsError() {
		return "", fmt.Errorf("Slack webhook returned status code %d: %s", resp.StatusCode(), resp.String())
	}

	return "", nil
}
