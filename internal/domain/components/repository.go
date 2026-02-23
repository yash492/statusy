package components

import (
	"context"
	"time"

	"github.com/yash492/statusy/internal/common"
)

type ComponentParams struct {
	Name             string
	ProviderID       string
	ServiceID        uint
	ComponentGroupID common.Nullable[uint]
}

type ComponentResult struct {
	ID               uint
	Name             string
	ProviderID       string
	ServiceID        uint
	ComponentGroupID common.Nullable[uint]
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        common.Nullable[time.Time]
}

type Repository interface {
	SaveAll(ctx context.Context, params []ComponentParams) ([]ComponentResult, error)
}
