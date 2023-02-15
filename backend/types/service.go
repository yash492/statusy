package types

type Service struct {
	BaseModel
	Name                   string `yaml:"name"`
	Slug                   string `yaml:"slug"`
	Link                   string `yaml:"link"`
	ShouldScrapWebsite     bool   `yaml:"should_scrap_website"`
	IncidentUrl            string `yaml:"incident_url"`
	ScheduleMaintenanceUrl string `yaml:"schedule_maintenance_url"`
	ComponentsUrl          string `yaml:"components_url"`
	ProviderType           string `yaml:"provider_type"`
}
