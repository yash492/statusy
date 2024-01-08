package dispatcher

import (
	"fmt"
	"log/slog"
	"sync"

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

	workerEvent, err := fetchSubscriptionContext(incident.IncidentUpdate.ID)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	for workerName, eventMap := range workerMap {
		slog.Info("dispatching event", "worker_name", workerName, "event_type", incident.State)
		worker := eventMap[incident.State]
		worker.Do(workerEvent)
	}
}

func fetchSubscriptionContext(incidentID uint) (types.WorkerEvent, error) {
	subscriptions, err := domain.Subscription.GetForIncidentUpdates(incidentID)
	if err != nil {
		slog.Error(err.Error())
		return types.WorkerEvent{}, err
	}

	if len(subscriptions) < 1 {
		slog.Error("no subscriptions found")
		return types.WorkerEvent{}, fmt.Errorf("no subscriptions were found for %v", incidentID)
	}

	components := lo.Map(subscriptions, func(subscription schema.SubscriptionForIncidentUpdates, _ int) types.ComponentsForWorker {
		return types.ComponentsForWorker{
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
		IncidentImpact:               subscriptions[0].IncidentImpact,
		IncidentUpdate:               subscriptions[0].IncidentUpdate,
		IncidentUpdateProviderStatus: subscriptions[0].IncidentUpdateProviderStatus,
		IncidentUpdateStatus:         subscriptions[0].IncidentUpdateStatus,
	}

	return workerEvent, nil
}
