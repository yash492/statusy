package services

import (
	"context"

	"github.com/yash492/statusy/internal/domain/statuspage"
)

type ServiceResult struct {
	ID                      uint
	Name                    string
	Slug                    string
	IncidentsUrl            string
	ScheduleMaintenancesUrl string
	ComponentsUrl           string
	ProviderType            statuspage.ProviderType
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
	SaveAll(ctx context.Context, services []ServiceParams) error
	GetAll(ctx context.Context) ([]ServiceResult, error)
}
