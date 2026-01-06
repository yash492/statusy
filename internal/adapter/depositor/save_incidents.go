package depositor

import (
	"github.com/yash492/statusy/internal/domain/statuspage"
)

// Assumptions:
// Incident Updates array will be sorted by time and will be of statuspage responsibility
func SaveIncidents(incidents []statuspage.Incident, components []statuspage.AggregateComponents) {
	// for _, incident := range incidents {
		// check whether incident exists in DB
		// if no then update the whole incident
		// else fetch the last incident update and compare
		// if the last incident update is same - do nothing
		// else add new incident update
		// lastIncidentUpdate := incident.Updates[len(incident.Updates)-1]
	// }
}
