package squadcast

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/jackc/pgx"
	"github.com/yash492/statusy/cmd/dispatcher/helpers"
	"github.com/yash492/statusy/pkg/domain"
	"github.com/yash492/statusy/pkg/types"
)

type IncidentEvent struct {
	Message     string       `json:"message"`
	Description string       `json:"description"`
	Tags        IncidentTags `json:"tags"`
	Priority    string       `json:"priority"`
	Status      string       `json:"status"`
	EventID     string       `json:"event_id"`
}

type IncidentTags map[string]Tag

type Tag struct {
	Value string `json:"value"`
	Color string `json:"color,omitempty"`
}

func (i IncidentOpenWorker) Do(event types.WorkerEvent) error {

	squadcast, err := domain.SquadcastExtension.Get()
	if err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return nil
		}
		return err
	}

	webhook := squadcast.WebhookURL

	client := resty.New()
	affectedComponents := helpers.ConvertComponentsToStr(event.Components)
	description := fmt.Sprintf("Created At: %v\n", event.IncidentUpdateStatusTime.UTC().Format(time.RFC822)) +
		fmt.Sprintf("Affected Components: %v\n", affectedComponents) +
		fmt.Sprintf("Incident Link: %v\n", event.IncidentLink) +
		fmt.Sprintln(event.IncidentUpdate)

	incidentEvent := IncidentEvent{
		Message:     fmt.Sprintf("%s: %s", event.ServiceName, event.IncidentName),
		Description: description,
		Tags: IncidentTags{
			"Service": Tag{
				Value: "Plivo",
			},
			"Components": Tag{
				Value: affectedComponents,
			},
			"Link": Tag{
				Value: event.IncidentLink,
			},
		},
		Status:  "trigger",
		EventID: makeEventID(event.ServiceID, event.IncidentID),
	}

	resp, err := client.R().SetBody(incidentEvent).Post(webhook)
	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("squadcast webhook failed with status code: %v and error: %v", resp.StatusCode(), resp.Error())
	}
	return nil
}
