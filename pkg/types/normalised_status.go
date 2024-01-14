package types

import "fmt"

const (
	Incident           string = "incident"
	IncidentTriggered  string = "triggered"
	IncidentInProgress string = "in_progress"
	IncidentResolved   string = "resolved"
)

var (
	IncidentTriggeredEventType  string = fmt.Sprintf("%v.%v", Incident, IncidentTriggered)
	IncidentInProgressEventType string = fmt.Sprintf("%v.%v", Incident, IncidentInProgress)
	IncidentResolvedEventType   string = fmt.Sprintf("%v.%v", Incident, IncidentResolved)
)
