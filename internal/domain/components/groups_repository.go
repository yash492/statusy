package components

import (
	"context"
	"time"

	"github.com/yash492/statusy/internal/common"
)

type GroupParams struct {
	Name       string
	ProviderID string
	ServiceID  uint
}

type ComponentGroupResult struct {
	ID         uint
	Name       string
	ProviderID string
	ServiceID  uint
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  common.Nullable[time.Time]
}

type GroupRepository interface {
	SaveAll(ctx context.Context, params []GroupParams) ([]ComponentGroupResult, error)
}
