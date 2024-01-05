package schema

type Service struct {
	BaseModel
	Name                   string `db:"name"`
	Link                   string `db:"link"`
	Slug                   string `db:"slug"`
	ProviderType           string `db:"provider_type"`
	ShouldScrapWebsite     bool   `db:"should_scrap_website"`
	IncidentURL            string `db:"incident_url"`
	ScheduleMaintenanceURL string `db:"schedule_maintenance_url"`
	ComponentsURL          string `db:"components_url"`
}
