package worker

import (
	"github.com/yash492/statusy/pkg/types"
)

type WorkerMap map[string]map[string]Worker

var dispatchWorker WorkerMap

func init() {
	dispatchWorker = make(WorkerMap, 0)
}

type Worker interface {
	Do(event types.WorkerEvent) error
}
