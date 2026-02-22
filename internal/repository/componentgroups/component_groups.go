package componentgroups

import (
	"context"
	"time"
)

type ComponentsGroupParams struct {
	Name       string
	ProviderID string
	ServiceID  uint
}

type ComponentsGroupResult struct {
	ID         uint       `db:"id"`
	Name       string     `db:"name"`
	ProviderID string     `db:"provider_id"`
	ServiceID  uint       `db:"service_id"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
	DeletedAt  *time.Time `db:"deleted_at"`
}

type ComponentsGroupRepository interface {
	SaveAll(ctx context.Context, params []ComponentsGroupParams) ([]ComponentsGroupResult, error)
}
