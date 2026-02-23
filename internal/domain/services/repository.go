package services

import (
	"context"
	"time"

	"github.com/yash492/statusy/internal/common"
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
	DeletedAt               common.Nullable[time.Time]
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
	GetAll(ctx context.Context) ([]ServiceResult, error)
}
