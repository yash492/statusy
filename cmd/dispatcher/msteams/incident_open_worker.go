package msteams

import (
	"github.com/yash492/statusy/pkg/types"
)

func (i IncidentOpenWorker) Do(event types.WorkerEvent) error {
	return dispatchMsTeamsMsg(event)
}
