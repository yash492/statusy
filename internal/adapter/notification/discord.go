package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/yash492/statusy/internal/common/jsonutil"
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
	data AlertData,
	prevExtID string,
) (string, error) {
	cfg, err := jsonutil.UnmarshalWithType[DiscordConfig](configRaw)
	if err != nil {
		return "", err
	}
	if cfg.WebhookURL == "" {
		return "", fmt.Errorf("Discord webhook URL is empty")
	}

	_, colorInt := getColor(data.Status)
	comps := formatComponents(data.Components)

	fields := []*discordgo.MessageEmbedField{
		{Name: "Service", Value: data.ServiceName, Inline: true},
		{Name: "Status", Value: data.Status.ForDisplay(), Inline: true},
		{Name: "Updated At", Value: data.UpdatedAt.UTC().Format(time.RFC822), Inline: true},
	}

	if data.StartTime != nil {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Start Time",
			Value:  data.StartTime.UTC().Format(time.RFC822),
			Inline: true,
		})
	}
	if data.EndTime != nil {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "End Time",
			Value:  data.EndTime.UTC().Format(time.RFC822),
			Inline: true,
		})
	}

	fields = append(fields,
		&discordgo.MessageEmbedField{Name: "Affected Components", Value: comps},
		&discordgo.MessageEmbedField{Name: "Description", Value: data.Description},
	)

	embed := &discordgo.MessageEmbed{
		Title:  fmt.Sprintf(":rotating_light: **%s - %s**", data.ServiceName, data.Title),
		URL:    data.Link,
		Color:  colorInt,
		Fields: fields,
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
