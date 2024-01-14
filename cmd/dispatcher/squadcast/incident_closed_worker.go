package squadcast

import (
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/jackc/pgx"
	"github.com/yash492/statusy/pkg/domain"
	"github.com/yash492/statusy/pkg/types"
)

func (i IncidentClosedWorker) Do(event types.WorkerEvent) error {
	squadcast, err := domain.SquadcastExtension.Get()
	if err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return nil
		}
		return err
	}

	webhook := squadcast.WebhookURL

	client := resty.New()
	resp, err := client.
		R().
		SetBody(
			types.JSON{
				"status":   "resolve",
				"event_id": makeEventID(event.ServiceID, event.IncidentID),
			}).
		Post(webhook)

	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("squadcast webhook failed with status code: %v and error: %v", resp.StatusCode(), resp.Error())
	}

	return nil
}
