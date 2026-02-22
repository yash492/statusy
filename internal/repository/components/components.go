package components

import (
	"context"
	"time"
)

type ComponentParams struct {
	Name             string `db:"name"`
	ProviderID       string `db:"provider_id"`
	ServiceID        uint   `db:"service_id"`
	ComponentGroupID *uint  `db:"component_group_id"`
}
type ComponentResult struct {
	ID               uint       `db:"id"`
	Name             string     `db:"name"`
	ProviderID       string     `db:"provider_id"`
	ServiceID        uint       `db:"service_id"`
	ComponentGroupID *uint      `db:"component_group_id"`
	CreatedAt        time.Time  `db:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at"`
	DeletedAt        *time.Time `db:"deleted_at"`
}

type ComponentRepository interface {
	SaveAll(ctx context.Context, params []ComponentParams) ([]ComponentResult, error)
}
