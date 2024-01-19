package slack

import (
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx"
	"github.com/slack-go/slack"
	"github.com/yash492/statusy/cmd/dispatcher/helpers"
	"github.com/yash492/statusy/pkg/domain"
	"github.com/yash492/statusy/pkg/types"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type IncidentOpenWorker struct{}
type IncidentInProgressWorker struct{}
type IncidentClosedWorker struct{}

var slackMsgColor = map[string]string{
	types.IncidentTriggeredEventType:  "#FF5733",
	types.IncidentInProgressEventType: "#FDDA0D",
	types.IncidentResolvedEventType:   "#4CBB17",
}

func dispatchSlackMsg(event types.WorkerEvent) error {
	slack, err := domain.ChatopsExtension.GetByType("slack")
	if err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return nil
		}
		return err
	}

	slackWebhookURL := slack.WebhookURL
	components := helpers.ConvertComponentsToStr(event.Components)

	blocks, attachment := makeSlackWebhookMessage(components, event)
	return sendWebhookMsg(slackWebhookURL, blocks, attachment)
}

func makeSlackWebhookMessage(components string, event types.WorkerEvent) (slack.Message, slack.Attachment) {
	msg := slack.NewBlockMessage(
		slack.SectionBlock{
			Type: slack.MBTSection,
			Text: &slack.TextBlockObject{
				Type: slack.MarkdownType,
				Text: fmt.Sprintf(":rotating_light: *%v*", helpers.SlackHyperlinkFormat(event.IncidentLink, event.IncidentName)),
			},
			Fields: []*slack.TextBlockObject{
				{
					Type: slack.MarkdownType,
					Text: fmt.Sprintf("*Service:* %v", event.ServiceName),
				},
				{
					Type: slack.MarkdownType,
					Text: fmt.Sprintf("*Status:* `%v`", cases.Title(language.AmericanEnglish).String(event.IncidentUpdateProviderStatus)),
				},
				{
					Type: slack.MarkdownType,
					Text: fmt.Sprintf("*Impact:* %v", cases.Title(language.AmericanEnglish).String(event.IncidentImpact)),
				},
				{
					Type: slack.MarkdownType,
					Text: fmt.Sprintf("*Created At:* `%v`", event.IncidentUpdateStatusTime.UTC().Format(time.RFC822)),
				},
				{
					Type: slack.MarkdownType,
					Text: fmt.Sprintf("*Affected Components:*\n %v", components),
				},
			},
		},
	)

	attachmentBlocks := slack.NewBlockMessage(slack.SectionBlock{
		Type: slack.MBTSection,
		Text: &slack.TextBlockObject{
			Type: slack.MarkdownType,
			Text: event.IncidentUpdate,
		},
	})

	attachment := slack.Attachment{
		Color:  slackMsgColor[event.EventType],
		Blocks: attachmentBlocks.Blocks,
	}

	return msg, attachment
}

func sendWebhookMsg(webhookURL string, blocks slack.Message, attachment slack.Attachment) error {
	err := slack.PostWebhook(webhookURL, &slack.WebhookMessage{
		Blocks:      &blocks.Blocks,
		Attachments: []slack.Attachment{attachment},
	})

	return err
}
