package types

import "time"

type WorkerEvent struct {
	ServiceID                    uint
	ServiceName                  string
	IncidentID                   uint
	IncidentName                 string
	IncidentLink                 string
	IncidentImpact               string
	IncidentUpdate               string
	IncidentUpdateID             uint
	IncidentUpdateProviderStatus string
	IncidentUpdateStatus         string
	Components                   []ComponentsWithNameAndID
	IsAllComponents              bool
	EventType                    string
	IncidentUpdateStatusTime     time.Time
}

type ComponentsWithNameAndID struct {
	Name string
	ID   uint
}
