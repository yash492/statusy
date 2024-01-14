package msteams

import "github.com/yash492/statusy/pkg/types"

func (i IncidentClosedWorker) Do(event types.WorkerEvent) error {
	return dispatchMsTeamsMsg(event)
}
