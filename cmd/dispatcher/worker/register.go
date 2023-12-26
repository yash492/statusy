package worker

import (
	"github.com/yash492/statusy/cmd/dispatcher/slack"
	"github.com/yash492/statusy/pkg/types"
)

func New() WorkerMap {
	registerWorker(types.SlackWorker, types.IncidentOpenEventType, slack.IncidentOpenWorker{})
	registerWorker(types.SlackWorker, types.IncidentInProgressEventType, slack.IncidentInProgressWorker{})
	registerWorker(types.SlackWorker, types.IncidentClosesEventType, slack.IncidentClosedWorker{})
	return dispatchWorker
}

func registerWorker(workerName string, eventType string, worker Worker) {
	eventWorkerMap, ok := dispatchWorker[workerName]
	if !ok {
		eventWorkerMap = make(map[string]Worker, 0)
	}
	eventWorkerMap[eventType] = worker
	dispatchWorker[workerName] = eventWorkerMap
}
