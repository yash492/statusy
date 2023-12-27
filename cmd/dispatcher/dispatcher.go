package dispatcher

import (
	"log/slog"
	"sync"

	"github.com/yash492/statusy/cmd/dispatcher/worker"
	"github.com/yash492/statusy/pkg/queue"
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
	for workerName, eventMap := range workerMap {
		slog.Info("dispatching event", "worker_name", workerName, "event_type", incident.State)
		worker := eventMap[incident.State]
		worker.Do(types.WorkerEvent{
			IncidentUpdate: incident.IncidentUpdate,
		})
	}
}
