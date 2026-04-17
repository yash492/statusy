package scheduledmaintenance

import (
	"context"
	"time"

	"github.com/yash492/statusy/internal/common/nullable"
)

type ScheduledMaintenanceParams struct {
	Title             string
	Link              string
	ProviderImpact    nullable.Nullable[string]
	Impact            nullable.Nullable[string]
	StartsAt          time.Time
	EndsAt            time.Time
	ServiceID         uint
	ProviderID        string
	ProviderCreatedAt time.Time
}

type ScheduledMaintenanceResult struct {
	ID                uint
	Title             string
	Link              string
	ProviderImpact    nullable.Nullable[string]
	Impact            nullable.Nullable[string]
	StartsAt          time.Time
	EndsAt            time.Time
	ServiceID         uint
	ProviderID        string
	ProviderCreatedAt time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         nullable.Nullable[time.Time]
}

type ScheduledMaintenanceByServiceParams struct {
	ServiceID uint
	Limit     int
	Offset    int
}

type ScheduledMaintenanceByServiceResult struct {
	ID                uint
	ServiceID         uint
	Title             string
	Status            string
	StartsAt          time.Time
	EndsAt            time.Time
	Link              string
	ProviderCreatedAt time.Time
}

type FeedScheduledMaintenanceByServiceResult struct {
	ID                 uint
	ServiceID          uint
	Title              string
	Status             string
	Link               string
	ProviderCreatedAt  time.Time
	AffectedComponents string
}

type ScheduledMaintenanceUpdateParams struct {
	ScheduledMaintenanceID uint
	Description            string
	ProviderID             string
	ProviderStatus         string
	Status                 string
	StatusTime             time.Time
}

type ScheduledMaintenanceUpdateResult struct {
	ID                     uint
	ScheduledMaintenanceID uint
	Description            string
	ProviderID             string
	ProviderStatus         string
	Status                 string
	StatusTime             time.Time
	CreatedAt              time.Time
	UpdatedAt              time.Time
	DeletedAt              nullable.Nullable[time.Time]
}

type ScheduledMaintenanceComponentParams struct {
	ScheduledMaintenanceID uint
	ComponentID            uint
}

type ScheduledMaintenanceComponentResult struct {
	ID                     uint
	ScheduledMaintenanceID uint
	ComponentID            uint
	CreatedAt              time.Time
	UpdatedAt              time.Time
	DeletedAt              nullable.Nullable[time.Time]
}

type Repository interface {
	SaveAll(ctx context.Context, params []ScheduledMaintenanceParams) ([]ScheduledMaintenanceResult, error)
	GetByService(ctx context.Context, params ScheduledMaintenanceByServiceParams) ([]ScheduledMaintenanceByServiceResult, error)
	GetFeedByService(ctx context.Context, params ScheduledMaintenanceByServiceParams) ([]FeedScheduledMaintenanceByServiceResult, error)
}

type UpdatesRepository interface {
	SaveAll(ctx context.Context, params []ScheduledMaintenanceUpdateParams) ([]ScheduledMaintenanceUpdateResult, error)
}

type ComponentsRepository interface {
	SaveAll(ctx context.Context, params []ScheduledMaintenanceComponentParams) ([]ScheduledMaintenanceComponentResult, error)
}
