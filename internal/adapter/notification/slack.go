package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/slack-go/slack"
	"github.com/yash492/statusy/internal/common/jsonutil"
)

type SlackConfig struct {
	WebhookURL string `json:"webhook_url"`
}

type SlackDispatcher struct{}

func newSlackDispatcher() *SlackDispatcher {
	return &SlackDispatcher{}
}

// Ensure SlackDispatcher implements ChannelDispatcher interface
var _ ChannelDispatcher = &SlackDispatcher{}

func (s *SlackDispatcher) Send(
	ctx context.Context,
	configRaw json.RawMessage,
	data AlertData,
	prevExtID string,
) (string, error) {
	cfg, err := jsonutil.UnmarshalWithType[SlackConfig](configRaw)
	if err != nil {
		return "", err
	}
	comps := formatComponents(data.Components)

	titleText := fmt.Sprintf(":rotating_light: *<%s|%s - %s>*", data.Link, data.ServiceName, data.Title)
	if data.Link == "" {
		titleText = fmt.Sprintf(":rotating_light: *%s - %s*", data.ServiceName, data.Title)
	}

	fieldsText := fmt.Sprintf(
		"*Service:* %s\n*Status:* `%s`\n*Updated At:* `%s`\n*Affected Components:*\n%s",
		data.ServiceName,
		data.Status.ForDisplay(),
		data.UpdatedAt.UTC().Format(time.RFC822),
		comps,
	)
	if data.StartTime != nil {
		fieldsText += fmt.Sprintf("\n*Start Time:* `%s`", data.StartTime.UTC().Format(time.RFC822))
	}
	if data.EndTime != nil {
		fieldsText += fmt.Sprintf("\n*End Time:* `%s`", data.EndTime.UTC().Format(time.RFC822))
	}

	msg := slack.NewBlockMessage(
		slack.NewSectionBlock(&slack.TextBlockObject{Type: slack.MarkdownType, Text: titleText}, nil, nil),
		slack.NewSectionBlock(&slack.TextBlockObject{Type: slack.MarkdownType, Text: fieldsText}, nil, nil),
		slack.NewSectionBlock(&slack.TextBlockObject{Type: slack.MarkdownType, Text: fmt.Sprintf("*Description:*\n%s", data.Description)}, nil, nil),
	)

	payload := &slack.WebhookMessage{
		Blocks: &msg.Blocks,
	}

	if err := slack.PostWebhookContext(ctx, cfg.WebhookURL, payload); err != nil {
		return "", fmt.Errorf("failed to send Slack message: %w", err)
	}

	return "", nil
}
