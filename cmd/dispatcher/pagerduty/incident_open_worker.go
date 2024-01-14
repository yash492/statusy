package pagerduty

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/PagerDuty/go-pagerduty"

	"github.com/jackc/pgx"
	"github.com/yash492/statusy/cmd/dispatcher/helpers"
	"github.com/yash492/statusy/pkg/domain"
	"github.com/yash492/statusy/pkg/types"
)

func (i IncidentOpenWorker) Do(event types.WorkerEvent) error {

	pagerdutyExtension, err := domain.PagerdutyExtension.Get()
	if err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return nil
		}
		return err
	}

	ctx := context.Background()
	routingKey := pagerdutyExtension.RoutingKey

	_, err = pagerduty.ManageEventWithContext(ctx, pagerduty.V2Event{
		RoutingKey: routingKey,
		Action:     "trigger",
		Client:     event.ServiceName,
		DedupKey:   makeEventID(event.ServiceID, event.IncidentID),
		ClientURL:  event.IncidentLink,
		Payload: &pagerduty.V2Payload{
			Timestamp: event.IncidentUpdateStatusTime.UTC().Format(time.RFC3339),
			Summary:   fmt.Sprintf("%v: %v", event.ServiceName, event.IncidentName),
			Source:    "Statusy",
			Severity:  "warning",
			Component: helpers.ConvertComponentsToStr(event.Components),
			Details: types.JSON{
				"Description": event.IncidentUpdate,
			},
		},
	})

	return err
}
