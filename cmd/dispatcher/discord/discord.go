package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/go-resty/resty/v2"
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
		return err
	}

	webhook := discord.WebhookURL

	result := types.JSON{}
	client := resty.New()

	discordMsg := discordgo.Message{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title: fmt.Sprintf(":rotating_light: **%v**", event.IncidentName),
				URL: event.IncidentLink,
				Color: msgColor[event.EventType],
				Fields: []*discordgo.MessageEmbedField{
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
					{
						Name:  "Affected Components",
						Value: helpers.ConvertComponentsToStr(event.Components),
					},
					{
						Name:  "Description",
						Value: event.IncidentUpdate,
					},
				},
			},
		},
	}

	resp, err := client.R().SetBody(discordMsg).SetResult(result).Post(webhook)
	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("discord webhook failed with status code: %v and error: %v", resp.StatusCode(), resp.Result())
	}

	return nil

}
