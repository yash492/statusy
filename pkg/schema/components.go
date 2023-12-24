package schema

type Component struct {
	BaseModel
	Name       string `db:"name"`
	ServiceID  uint   `db:"service_id"`
	ProviderID string `db:"provider_id"`
}
