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
}

type ComponentsForWorker struct {
	Name string
	ID   uint
}
