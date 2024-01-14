package slack

import (
	"errors"
	"fmt"

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

	blocks, attachment := makeSlackWebhookMessage(components, event.IncidentName, event.ServiceName, event.IncidentUpdateProviderStatus, event.EventType, event.IncidentUpdate, event.IncidentLink)
	return sendWebhookMsg(slackWebhookURL, blocks, attachment)
}

func makeSlackWebhookMessage(components string, incidentName, serviceName, providerStatus, eventType, incidentUpdateDescription, incidentLink string) (slack.Message, slack.Attachment) {
	msg := slack.NewBlockMessage(
		slack.SectionBlock{
			Type: slack.MBTSection,
			Text: &slack.TextBlockObject{
				Type: slack.MarkdownType,
				Text: fmt.Sprintf(":rotating_light: *%v*", helpers.SlackHyperlinkFormat(incidentLink, incidentName)),
			},
			Fields: []*slack.TextBlockObject{
				{
					Type: slack.MarkdownType,
					Text: fmt.Sprintf("*Service:* %v", serviceName),
				},
				{
					Type: slack.MarkdownType,
					Text: fmt.Sprintf("*Status:* `%v`", cases.Title(language.AmericanEnglish).String(providerStatus)),
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
			Text: incidentUpdateDescription,
		},
	})

	attachment := slack.Attachment{
		Color:  slackMsgColor[eventType],
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
