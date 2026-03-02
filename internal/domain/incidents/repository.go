package incidents

import (
	"context"
	"time"

	"github.com/yash492/statusy/internal/common/nullable"
)

type IncidentParams struct {
	Name              string
	Link              string
	ProviderImpact    nullable.Nullable[string]
	Impact            nullable.Nullable[string]
	ServiceID         uint
	ProviderID        string
	ProviderCreatedAt time.Time
}

type IncidentResult struct {
	ID                uint
	Name              string
	Link              string
	ProviderImpact    nullable.Nullable[string]
	Impact            nullable.Nullable[string]
	ServiceID         uint
	ProviderID        string
	ProviderCreatedAt time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         nullable.Nullable[time.Time]
}

type IncidentUpdateParams struct {
	IncidentID     uint
	Description    string
	ProviderID     string
	ProviderStatus string
	Status         string
	StatusTime     time.Time
}

type IncidentUpdateResult struct {
	ID             uint
	IncidentID     uint
	Description    string
	ProviderID     string
	ProviderStatus string
	Status         string
	StatusTime     time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      nullable.Nullable[time.Time]
}

type IncidentComponentParams struct {
	IncidentID  uint
	ComponentID uint
}

type IncidentComponentResult struct {
	ID          uint
	IncidentID  uint
	ComponentID uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   nullable.Nullable[time.Time]
}

type Repository interface {
	SaveAll(ctx context.Context, params []IncidentParams) ([]IncidentResult, error)
}

type UpdatesRepository interface {
	SaveAll(ctx context.Context, params []IncidentUpdateParams) ([]IncidentUpdateResult, error)
}

type ComponentsRepository interface {
	SaveAll(ctx context.Context, params []IncidentComponentParams) ([]IncidentComponentResult, error)
}
