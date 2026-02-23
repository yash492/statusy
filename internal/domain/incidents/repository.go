package incidents

import (
	"context"
	"time"

	"github.com/yash492/statusy/internal/common"
)

type IncidentParams struct {
	Name              string
	Link              string
	ProviderImpact    common.Nullable[string]
	Impact            common.Nullable[string]
	ServiceID         uint
	ProviderID        string
	ProviderCreatedAt time.Time
}

type IncidentResult struct {
	ID                uint
	Name              string
	Link              string
	ProviderImpact    common.Nullable[string]
	Impact            common.Nullable[string]
	ServiceID         uint
	ProviderID        string
	ProviderCreatedAt time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         common.Nullable[time.Time]
}

type Repository interface {
	SaveAll(ctx context.Context, params []IncidentParams) ([]IncidentResult, error)
}
