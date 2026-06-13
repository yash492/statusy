package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

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
	isFirst bool,
	isResolve bool,
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
	if isResolve {
		teamsColor = adaptivecard.ColorGood
	} else if isFirst {
		teamsColor = adaptivecard.ColorAttention
	} else {
		teamsColor = adaptivecard.ColorWarning
	}

	title := fmt.Sprintf("%s Alert: %s", data.ServiceName, data.Title)
	card, err := adaptivecard.NewTextBlockCard(data.Description, title, true)
	if err != nil {
		return "", fmt.Errorf("failed to create Teams adaptive card: %w", err)
	}

	statusBlock := adaptivecard.NewTextBlock(fmt.Sprintf("Status: %s", data.Status), true)
	statusBlock.Color = teamsColor
	statusBlock.Weight = "Bolder"
	_ = card.AddElement(false, statusBlock)

	comps := formatComponents(data.Components)
	factSet := adaptivecard.NewFactSet()
	_ = factSet.AddFact(
		adaptivecard.Fact{Title: "Status:", Value: data.Status},
		adaptivecard.Fact{Title: "Affected Components:", Value: comps},
		adaptivecard.Fact{Title: "Updated At:", Value: data.UpdatedAt.UTC().Format("2006-01-02 15:04:05 MST")},
	)
	_ = card.AddElement(false, adaptivecard.Element(factSet))

	if data.Link != "" {
		urlAction, err := adaptivecard.NewActionOpenURL(data.Link, "View Status Page")
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
