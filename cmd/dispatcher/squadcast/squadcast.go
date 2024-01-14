package squadcast

import "fmt"

type IncidentOpenWorker struct{}
type IncidentInProgressWorker struct{}
type IncidentClosedWorker struct{}

func makeEventID(serviceID, incidentID uint) string {
	return fmt.Sprintf("%d-%d", serviceID, incidentID)
}
