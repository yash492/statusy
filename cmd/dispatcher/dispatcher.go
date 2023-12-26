package dispatcher

import (
	"log/slog"
	"sync"

	"github.com/yash492/statusy/cmd/dispatcher/worker"
	"github.com/yash492/statusy/pkg/queue"
)

func New(queue *queue.Queue, wg *sync.WaitGroup) error {
	workerMap := worker.New()

	for  {
		incident, err := queue.Consume()
		if err != nil {
			continue
		}
		dispatchIncident(incident, workerMap)
	}

	wg.Done()
	return nil
}

func dispatchIncident(incident queue.IncidentPayload, workerMap worker.WorkerMap) {
	for workerName, eventMap := range workerMap {
		for eventType, worker := range eventMap {
			worker.Do()
			slog.Info("Dispatching worker for event type", "worker_name", workerName, "event_type", eventType)
		}
		slog.Info("Dispatching worker", "worker_name", workerName)

	}
}
