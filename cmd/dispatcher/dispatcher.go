package dispatcher

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/samber/lo"
	"github.com/yash492/statusy/cmd/dispatcher/worker"
	"github.com/yash492/statusy/pkg/domain"
	"github.com/yash492/statusy/pkg/queue"
	"github.com/yash492/statusy/pkg/schema"
	"github.com/yash492/statusy/pkg/types"
)

func New(queue *queue.Queue, wg *sync.WaitGroup) {
	workerMap := worker.New()

	for {
		incident, err := queue.Consume()
		if err != nil {
			continue
		}
		dispatchIncident(incident, workerMap)
	}

}

func dispatchIncident(incident queue.IncidentPayload, workerMap worker.WorkerMap) {

	workerEvent, err := fetchSubscriptionContext(incident.IncidentUpdate.ID, incident.State)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	wg := &sync.WaitGroup{}

	for workerName, eventMap := range workerMap {
		slog.Info("dispatching event", "worker_name", workerName, "event_type", incident.State, "incident_update_id", incident.IncidentUpdate.ID, "incident_id", workerEvent.IncidentID)
		worker, ok := eventMap[workerEvent.EventType]
		if !ok {
			return
		}

		wg.Add(1)
		go func(w types.WorkerEvent, wg *sync.WaitGroup) {
			defer wg.Done()
			err := worker.Do(w)
			if err != nil {
				slog.Error(err.Error())
			}
			// Limiting the request speed as not to cross rate limits
			// and giving ample breathing room for services to process the request
			time.Sleep(1 * time.Second)
		}(workerEvent, wg)

	}

	wg.Wait()
}

func fetchSubscriptionContext(incidentID uint, eventType string) (types.WorkerEvent, error) {
	subscriptions, err := domain.Subscription.GetForIncidentUpdates(incidentID)
	if err != nil {
		slog.Error(err.Error())
		return types.WorkerEvent{}, err
	}

	if len(subscriptions) < 1 {
		slog.Error("no subscriptions found")
		return types.WorkerEvent{}, fmt.Errorf("no subscriptions were found for %v", incidentID)
	}

	components := lo.Map(subscriptions, func(subscription schema.SubscriptionForIncidentUpdate, _ int) types.ComponentsWithNameAndID {
		return types.ComponentsWithNameAndID{
			Name: subscription.ComponentName,
			ID:   (subscription.ComponentID),
		}
	})

	workerEvent := types.WorkerEvent{
		Components:                   components,
		ServiceID:                    subscriptions[0].ServiceID,
		ServiceName:                  subscriptions[0].ServiceName,
		IncidentID:                   subscriptions[0].IncidentID,
		IncidentName:                 subscriptions[0].IncidentName,
		IncidentLink:                 subscriptions[0].IncidentLink,
		IncidentImpact:               subscriptions[0].IncidentImpact.String,
		IncidentUpdate:               subscriptions[0].IncidentUpdate,
		IncidentUpdateProviderStatus: subscriptions[0].IncidentUpdateProviderStatus,
		IncidentUpdateStatus:         subscriptions[0].IncidentUpdateStatus,
		IsAllComponents:              subscriptions[0].IsAllComponents,
		EventType:                    eventType,
		IncidentUpdateStatusTime:     subscriptions[0].IncidentUpdateStatusTime,
		IncidentUpdateID:             subscriptions[0].IncidentUpdateID,
	}

	return workerEvent, nil
}
