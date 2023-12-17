package schema

type tableName string

type Service struct {
	BaseModel
	Name                   string
	Link                   string
	Slug                   string
	ProviderType           string
	ShouldScrapWebsite     bool
	IncidentURL            string
	ScheduleMaintenanceURL string
	ComponentsURL          string
}

const ServiceTable tableName = "services"

func (t tableName) String() string {
	return string(t)
}

func ServiceColumns() []string {
	return []string{
		"id",
		"name",
		"slug",
		"link",
		"shoud_scrap_website",
		"incident_url",
		"schedule_maintenance_url",
		"components_url",
		"provider_type",
	}
}
