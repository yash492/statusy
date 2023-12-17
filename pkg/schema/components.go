package schema

type Component struct {
	BaseModel
	Name                string `db:"name"`
	ServiceId           uint   `db:"service_id"`
	ProviderComponentId string `db:"provider_component_id"`
}
