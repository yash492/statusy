package components

import (
	"context"
	"time"

	"github.com/yash492/statusy/internal/common/nullable"
)

type ComponentParams struct {
	Name             string
	ProviderID       string
	ServiceID        uint
	ComponentGroupID nullable.Nullable[uint]
}

type ComponentResult struct {
	ID               uint
	Name             string
	ProviderID       string
	ServiceID        uint
	ComponentGroupID nullable.Nullable[uint]
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        nullable.Nullable[time.Time]
}

type Repository interface {
	SaveAll(ctx context.Context, params []ComponentParams) ([]ComponentResult, error)
}
