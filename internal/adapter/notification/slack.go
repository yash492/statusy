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
	Token      string `json:"token"`
	ChannelID  string `json:"channel_id"`
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
	if cfg.Token == "" {
		return "", fmt.Errorf("Slack token is empty")
	}
	if cfg.ChannelID == "" {
		return "", fmt.Errorf("Slack channel ID is empty")
	}

	colorHex, _ := getColor(data.Status)
	comps := formatComponents(data.Components)

	titleText := fmt.Sprintf(":rotating_light: *<%s|%s - %s>*", data.Link, data.ServiceName, data.Title)
	if data.Link == "" {
		titleText = fmt.Sprintf(":rotating_light: *%s - %s*", data.ServiceName, data.Title)
	}

	fields := []*slack.TextBlockObject{
		{Type: slack.MarkdownType, Text: fmt.Sprintf("*Service:* %s", data.ServiceName)},
		{Type: slack.MarkdownType, Text: fmt.Sprintf("*Status:* `%s`", data.Status.ForDisplay())},
		{Type: slack.MarkdownType, Text: fmt.Sprintf("*Updated At:* `%s`", data.UpdatedAt.UTC().Format(time.RFC822))},
		{Type: slack.MarkdownType, Text: fmt.Sprintf("*Affected Components:*\n%s", comps)},
	}

	if data.StartTime != nil {
		fields = append(fields, &slack.TextBlockObject{
			Type: slack.MarkdownType,
			Text: fmt.Sprintf("*Start Time:* `%s`", data.StartTime.UTC().Format(time.RFC822)),
		})
	}
	if data.EndTime != nil {
		fields = append(fields, &slack.TextBlockObject{
			Type: slack.MarkdownType,
			Text: fmt.Sprintf("*End Time:* `%s`", data.EndTime.UTC().Format(time.RFC822)),
		})
	}

	msg := slack.NewBlockMessage(
		slack.SectionBlock{
			Type: slack.MBTSection,
			Text: &slack.TextBlockObject{
				Type: slack.MarkdownType,
				Text: titleText,
			},
			Fields: fields,
		},
	)

	attachmentBlocks := slack.NewBlockMessage(slack.SectionBlock{
		Type: slack.MBTSection,
		Text: &slack.TextBlockObject{
			Type: slack.MarkdownType,
			Text: data.Description,
		},
	})

	attachment := slack.Attachment{
		Color:  colorHex,
		Blocks: attachmentBlocks.Blocks,
	}

	payload := &slack.WebhookMessage{
		Blocks:      &msg.Blocks,
		Attachments: []slack.Attachment{attachment},
	}

	if err := slack.PostWebhookContext(ctx, cfg.WebhookURL, payload); err != nil {
		return "", fmt.Errorf("failed to send Slack message: %w", err)
	}

	return "", nil
}
