package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"resty.dev/v3"
)

type DiscordConfig struct {
	WebhookURL string `json:"webhook_url"`
}

type DiscordDispatcher struct {
	client *resty.Client
	lg     *slog.Logger
}

func newDiscordDispatcher(client *resty.Client, lg *slog.Logger) *DiscordDispatcher {
	return &DiscordDispatcher{client: client, lg: lg}
}

// Ensure DiscordDispatcher implements ChannelDispatcher interface
var _ ChannelDispatcher = &DiscordDispatcher{}

func (d *DiscordDispatcher) Send(
	ctx context.Context,
	configRaw json.RawMessage,
	isFirst bool,
	isResolve bool,
	data AlertData,
	prevExtID string,
) (string, error) {
	var cfg DiscordConfig
	if err := json.Unmarshal(configRaw, &cfg); err != nil {
		return "", fmt.Errorf("failed to parse Discord config: %w", err)
	}
	if cfg.WebhookURL == "" {
		return "", fmt.Errorf("Discord webhook URL is empty")
	}

	_, colorInt := getColor(isFirst, isResolve)
	comps := formatComponents(data.Components)

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("[%s] Alert: %s", data.ServiceName, data.Title),
		Description: data.Description,
		URL:         data.Link,
		Color:       colorInt,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Status",
				Value:  data.Status,
				Inline: true,
			},
			{
				Name:   "Affected Components",
				Value:  comps,
				Inline: true,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Statusy",
		},
		Timestamp: data.UpdatedAt.UTC().Format(time.RFC3339),
	}

	payload := discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{embed},
	}

	if prevExtID != "" {
		patchURL := fmt.Sprintf("%s/messages/%s", strings.TrimSuffix(cfg.WebhookURL, "/"), prevExtID)
		resp, err := d.client.R().
			SetContext(ctx).
			SetBody(payload).
			Patch(patchURL)
		if err == nil && !resp.IsError() {
			return prevExtID, nil
		}

		if resp != nil && resp.StatusCode() == 404 {
			d.lg.WarnContext(ctx, "Discord message to update was not found (404), falling back to POST", slog.String("message_id", prevExtID))
		} else if err != nil {
			d.lg.WarnContext(ctx, "Discord PATCH failed, falling back to POST", slog.Any("err", err))
		} else {
			d.lg.WarnContext(ctx, "Discord PATCH returned error, falling back to POST", slog.Int("code", resp.StatusCode()), slog.String("body", resp.String()))
		}
	}

	postURL := cfg.WebhookURL
	if !strings.Contains(postURL, "?") {
		postURL += "?wait=true"
	} else {
		postURL += "&wait=true"
	}

	resp, err := d.client.R().
		SetContext(ctx).
		SetBody(payload).
		Post(postURL)
	if err != nil {
		return "", fmt.Errorf("failed to send Discord POST request: %w", err)
	}
	if resp.IsError() {
		return "", fmt.Errorf("Discord webhook returned status code %d: %s", resp.StatusCode(), resp.String())
	}

	var discordResp struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(resp.Bytes(), &discordResp); err != nil {
		d.lg.WarnContext(ctx, "failed to parse Discord message response ID", slog.Any("err", err))
		return "", nil
	}

	return discordResp.ID, nil
}
