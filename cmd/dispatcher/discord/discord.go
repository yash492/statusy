package discord

import (
	"errors"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-resty/resty/v2"
	"github.com/jackc/pgx"
	"github.com/yash492/statusy/cmd/dispatcher/helpers"
	"github.com/yash492/statusy/pkg/domain"
	"github.com/yash492/statusy/pkg/types"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type IncidentOpenWorker struct{}
type IncidentInProgressWorker struct{}
type IncidentClosedWorker struct{}

var msgColor = map[string]int{
	types.IncidentTriggeredEventType:  16729344,
	types.IncidentInProgressEventType: 16776960,
	types.IncidentResolvedEventType:   5763719,
}

func dispatchDiscordMsg(event types.WorkerEvent) error {

	discord, err := domain.ChatopsExtension.GetByType("discord")
	if err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return nil
		}
		return err
	}

	webhook := discord.WebhookURL

	result := types.JSON{}
	client := resty.New()

	fields := []*discordgo.MessageEmbedField{
		{
			Name:   "Service",
			Value:  event.ServiceName,
			Inline: true,
		},
		{
			Name:   "Incident Status",
			Value:  cases.Title(language.AmericanEnglish).String(event.IncidentUpdateProviderStatus),
			Inline: true,
		},
	}

	if event.IncidentImpact != "" {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Incident Impact",
			Value:  cases.Title(language.AmericanEnglish).String(event.IncidentImpact),
			Inline: true,
		})

	}

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Created At",
		Value:  event.IncidentUpdateStatusTime.UTC().Format(time.RFC822),
		Inline: true,
	},
		&discordgo.MessageEmbedField{
			Name:  "Affected Components",
			Value: helpers.ConvertComponentsToStr(event.Components),
		},
		&discordgo.MessageEmbedField{
			Name:  "Description",
			Value: event.IncidentUpdate,
		})

	discordMsg := discordgo.Message{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:  fmt.Sprintf(":rotating_light: **%v**", event.IncidentName),
				URL:    event.IncidentLink,
				Color:  msgColor[event.EventType],
				Fields: fields,
			},
		},
	}

	resp, err := client.R().SetBody(discordMsg).SetError(result).Post(webhook)
	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("discord webhook failed with status code: %v and error: %v", resp.StatusCode(), resp.Error())
	}

	return nil

}
