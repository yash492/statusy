package types

type WorkerEvent struct {
	ServiceID                    uint
	ServiceName                  string
	IncidentID                   uint
	IncidentName                 string
	IncidentLink                 string
	IncidentImpact               string
	IncidentUpdate               string
	IncidentUpdateProviderStatus string
	IncidentUpdateStatus         string
	Components                   []ComponentsForWorker
	IsAllComponents              bool
	EventType                    string
}

type ComponentsForWorker struct {
	Name string
	ID   uint
}
