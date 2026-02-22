package services

import (
	"context"
	"time"

	"github.com/yash492/statusy/internal/domain/statuspage"
)

type ServiceResult struct {
	ID                      uint                    `db:"id"`
	Name                    string                  `db:"name"`
	Slug                    string                  `db:"slug"`
	IncidentsUrl            string                  `db:"incidents_url"`
	ScheduleMaintenancesUrl string                  `db:"schedule_maintenances_url"`
	ComponentsUrl           string                  `db:"components_url"`
	ProviderType            statuspage.ProviderType `db:"provider_type"`
	CreatedAt               time.Time               `db:"created_at"`
	UpdatedAt               time.Time               `db:"updated_at"`
	DeletedAt               *time.Time              `db:"deleted_at"`
}

type ServiceParams struct {
	Name                    string                  `yaml:"name"`
	Slug                    string                  `yaml:"slug"`
	IncidentsUrl            string                  `yaml:"incidents_url"`
	ScheduleMaintenancesUrl string                  `yaml:"schedule_maintenances_url"`
	ComponentsUrl           string                  `yaml:"components_url"`
	ProviderType            statuspage.ProviderType `yaml:"provider_type"`
}

type ServiceRepository interface {
	SaveAll(ctx context.Context, services []ServiceParams) ([]ServiceResult, error)
	GetAll(ctx context.Context) ([]ServiceResult, error)
}
