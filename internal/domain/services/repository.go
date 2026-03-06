package services

import (
	"context"
	"time"

	"github.com/yash492/statusy/internal/common/nullable"
)

type ServiceResult struct {
	ID                      uint
	Name                    string
	Slug                    string
	IncidentsUrl            string
	ScheduleMaintenancesUrl string
	ComponentsUrl           string
	ProviderType            ProviderType
	CreatedAt               time.Time
	UpdatedAt               time.Time
	DeletedAt               nullable.Nullable[time.Time]
}

type ServiceParams struct {
	Name                    string       `yaml:"name"`
	Slug                    string       `yaml:"slug"`
	IncidentsUrl            string       `yaml:"incidents_url"`
	ScheduleMaintenancesUrl string       `yaml:"schedule_maintenances_url"`
	ComponentsUrl           string       `yaml:"components_url"`
	ProviderType            ProviderType `yaml:"provider_type"`
}

type Repository interface {
	SaveAll(ctx context.Context, services []ServiceParams) ([]ServiceResult, error)
	SearchBySlug(ctx context.Context, slug string) ([]ServiceResult, error)
	GetBySlug(ctx context.Context, slug string) (ServiceResult, error)
}
