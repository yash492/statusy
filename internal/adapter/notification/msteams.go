package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
	"github.com/atc0005/go-teams-notify/v2/adaptivecard"
	"github.com/yash492/statusy/internal/common/jsonutil"
)

type MsTeamsConfig struct {
	WebhookURL string `json:"webhook_url"`
}

type MsTeamsDispatcher struct {
	lg *slog.Logger
}

func newMsTeamsDispatcher(lg *slog.Logger) *MsTeamsDispatcher {
	return &MsTeamsDispatcher{lg: lg}
}

// Ensure MsTeamsDispatcher implements ChannelDispatcher interface
var _ ChannelDispatcher = &MsTeamsDispatcher{}

func (m *MsTeamsDispatcher) Send(
	ctx context.Context,
	configRaw json.RawMessage,
	data AlertData,
	prevExtID string,
) (string, error) {
	cfg, err := jsonutil.UnmarshalWithType[MsTeamsConfig](configRaw)
	if err != nil {
		return "", err
	}
	if cfg.WebhookURL == "" {
		return "", fmt.Errorf("MS Teams webhook URL is empty")
	}

	teamsColor := adaptivecard.ColorDefault
	if data.Status.IsResolved() {
		teamsColor = adaptivecard.ColorGood
	} else if data.Status.IsInitial() {
		teamsColor = adaptivecard.ColorAttention
	} else {
		teamsColor = adaptivecard.ColorWarning
	}

	comps := formatComponents(data.Components)

	card, err := adaptivecard.NewTextBlockCard(data.Description, fmt.Sprintf("🚨 %s - %s", data.ServiceName, data.Title), true)
	if err != nil {
		return "", fmt.Errorf("failed to create Teams adaptive card: %w", err)
	}

	statusBlock := adaptivecard.NewTextBlock(fmt.Sprintf("Status: %s", data.Status.ForDisplay()), true)
	statusBlock.Color = teamsColor
	statusBlock.Weight = "Bolder"
	_ = card.AddElement(false, statusBlock)

	facts := []adaptivecard.Fact{
		{Title: "Service:", Value: data.ServiceName},
		{Title: "Status:", Value: data.Status.ForDisplay()},
		{Title: "Affected Components:", Value: comps},
		{Title: "Updated At:", Value: data.UpdatedAt.UTC().Format(time.RFC822)},
	}

	if data.StartTime != nil {
		facts = append(facts, adaptivecard.Fact{Title: "Start Time:", Value: data.StartTime.UTC().Format(time.RFC822)})
	}
	if data.EndTime != nil {
		facts = append(facts, adaptivecard.Fact{Title: "End Time:", Value: data.EndTime.UTC().Format(time.RFC822)})
	}

	factSet := adaptivecard.NewFactSet()
	_ = factSet.AddFact(facts...)
	_ = card.AddElement(false, adaptivecard.Element(factSet))

	if data.Link != "" {
		urlAction, err := adaptivecard.NewActionOpenURL(data.Link, "View on Status Page")
		if err == nil {
			_ = card.AddAction(false, urlAction)
		}
	}

	msg, err := adaptivecard.NewMessageFromCard(card)
	if err != nil {
		return "", fmt.Errorf("failed to create Teams message from card: %w", err)
	}

	mstClient := goteamsnotify.NewTeamsClient()
	if err := mstClient.SendWithContext(ctx, cfg.WebhookURL, msg); err != nil {
		return "", fmt.Errorf("failed to send Teams message: %w", err)
	}

	return "", nil
}
