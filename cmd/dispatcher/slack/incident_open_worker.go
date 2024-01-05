package slack

import (
	"log/slog"

	"github.com/yash492/statusy/pkg/types"
)

func (i IncidentOpenWorker) Do(event types.WorkerEvent) error {
	slog.Info("slack incident open worker", "event_type", "incident.open", "status", event.IncidentUpdate.ProviderStatus)
	return nil
}
