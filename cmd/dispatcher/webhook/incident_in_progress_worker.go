package webhook

import "github.com/yash492/statusy/pkg/types"

func (i IncidentInProgressWorker) Do(event types.WorkerEvent) error {
	return dispatchWebhookMsg(event)
}
