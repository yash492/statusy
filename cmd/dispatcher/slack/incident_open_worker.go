package slack

import (
	"fmt"

	"github.com/slack-go/slack"
	"github.com/yash492/statusy/pkg/types"
)

func (i IncidentOpenWorker) Do(event types.WorkerEvent) error {

	slackWebhookURL := "https://hooks.slack.com/services/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"
	err := slack.PostWebhook(slackWebhookURL, &slack.WebhookMessage{
		Blocks: &slack.Blocks{
			BlockSet: []slack.Block{
				slack.SectionBlock{
					Type: slack.MBTSection,
					Text: &slack.TextBlockObject{
						Type: slack.MarkdownType,
						Text: fmt.Sprint(event.IncidentID) + " " + event.IncidentName,
					},
				},
				slack.SectionBlock{
					Type: slack.MBTSection,
					Text: &slack.TextBlockObject{
						Type: slack.PlainTextType,
						Text: fmt.Sprint(event.IncidentUpdate) + " " + event.IncidentUpdateProviderStatus,
					},
				},
			},
		},
	})

	if err != nil {
		return err
	}

	return nil
}
