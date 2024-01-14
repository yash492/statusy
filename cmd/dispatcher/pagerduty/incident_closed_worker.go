package pagerduty

import (
	"context"
	"errors"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/jackc/pgx"
	"github.com/yash492/statusy/pkg/domain"
	"github.com/yash492/statusy/pkg/types"
)

func (i IncidentClosedWorker) Do(event types.WorkerEvent) error {
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
		Action:     "resolve",
		DedupKey:   makeEventID(event.ServiceID, event.IncidentID),
	})

	return err
}
