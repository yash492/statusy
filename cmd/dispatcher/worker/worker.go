package worker

type WorkerMap map[string]map[string]Worker

var dispatchWorker WorkerMap

func init() {
	dispatchWorker = make(WorkerMap, 0)
}

type Worker interface {
	Do() error
}
