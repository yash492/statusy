package types

import "fmt"

const (
	Incident           string = "incident"
	IncidentOpen       string = "open"
	IncidentInProgress string = "in_progress"
	IncidentClosed     string = "closed"
)

var (
	IncidentOpenEventType       string = fmt.Sprintf("%v.%v", Incident, IncidentOpen)
	IncidentInProgressEventType string = fmt.Sprintf("%v.%v", Incident, IncidentInProgress)
	IncidentClosedEventType     string = fmt.Sprintf("%v.%v", Incident, IncidentClosed)
)
