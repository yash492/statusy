package msteams

import (
	"fmt"
	"strings"
	"time"

	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
	"github.com/atc0005/go-teams-notify/v2/messagecard"
	"github.com/yash492/statusy/cmd/dispatcher/helpers"
	"github.com/yash492/statusy/pkg/domain"
	"github.com/yash492/statusy/pkg/types"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type IncidentOpenWorker struct{}
type IncidentInProgressWorker struct{}
type IncidentClosedWorker struct{}

func dispatchMsTeamsMsg(event types.WorkerEvent) error {

	msteam, err := domain.ChatopsExtension.GetByType("msteams")
	if err != nil {
		return err
	}

	webhook := msteam.WebhookURL
	affectedComponents := helpers.ConvertComponentsToStr(event.Components)

	mstClient := goteamsnotify.NewTeamsClient()
	newLine := "<br/>"

	// Setup message card.
	msgCard := messagecard.NewMessageCard()
	msgCard.Title = fmt.Sprintf("🚨 %v", event.IncidentName)

	msgCardText := []string{
		fmt.Sprintf("**%v**", helpers.MarkdownHyperLinkFormat("Incident Link", event.IncidentLink)),
		newLine,
		fmt.Sprintf("**Service:** %v", event.ServiceName),
		newLine,
		fmt.Sprintf("**Incident Status:** %v", cases.Title(language.AmericanEnglish).String(event.IncidentUpdateProviderStatus)),
		newLine,
		"**Affected Components:**",
		newLine,
		affectedComponents,
		newLine,
		"**Created At**",
		newLine,
		event.IncidentUpdateStatusTime.UTC().Format(time.RFC822),
		newLine,
	}

	if event.IncidentImpact != "" {
		msgCardText = append(msgCardText,
			"**Impact**",
			newLine,
			cases.Title(language.AmericanEnglish).String(event.IncidentImpact),
			newLine,
		)
	}

	msgCardText = append(msgCardText,
		newLine,
		"**Description:**",
		newLine,
		event.IncidentUpdate,
	)

	msgCard.Text = strings.Join(msgCardText, "")

	// Send the message with default timeout/retry settings.
	return mstClient.Send(webhook, msgCard)

}
