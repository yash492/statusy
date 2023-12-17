package resource

import (
	"fmt"
	"log/slog"

	"github.com/gosimple/slug"
	"github.com/yash492/statusy/pkg/domain"
	"github.com/yash492/statusy/pkg/schema"
	"github.com/yash492/statusy/pkg/static"
	"gopkg.in/yaml.v3"
)

type services struct {
	Name                   string `yaml:"name"`
	Link                   string `yaml:"link"`
	ShouldScrapWebsite     bool   `yaml:"should_scrap_website"`
	IncidentURL            string `yaml:"incident_url"`
	ScheduleMaintenanceURL string `yaml:"schedule_maintenance_url"`
	ComponentsURL          string `yaml:"components_url"`
	ProviderType           string `yaml:"provider_type"`
}

func initServices() error {
	services, err := fetchServiceFromFile()
	if err != nil {
		return err
	}

	err = domain.Service.Create(services)
	if err != nil {
		return err
	}

	return nil
}

func fetchServiceFromFile() ([]schema.Service, error) {

	var servicesYml []services

	// Reading service yaml content in bytes
	bytes, err := static.Fs.ReadFile("services.yml")
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	err = yaml.Unmarshal(bytes, &servicesYml)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	services, err := addAndVerifySlugToServices(servicesYml)
	if err != nil {
		slog.Error("error while adding slug to services %s", err)
		return nil, err
	}

	return services, nil
}

func addAndVerifySlugToServices(servicesYml []services) ([]schema.Service, error) {
	var services []schema.Service
	var slugMap map[string]bool

	for _, service := range servicesYml {
		serviceSlug := slug.Make(service.Name)
		exists := slugMap[serviceSlug]
		if exists {
			return nil, fmt.Errorf("slug already exists for the %v service, please change the name", service.Name)
		} else {
			services = append(services, schema.Service{
				Name:                   service.Name,
				Slug:                   serviceSlug,
				Link:                   service.Link,
				ShouldScrapWebsite:     service.ShouldScrapWebsite,
				IncidentURL:            service.IncidentURL,
				ScheduleMaintenanceURL: service.ScheduleMaintenanceURL,
				ComponentsURL:          service.ComponentsURL,
				ProviderType:           service.ProviderType,
			})
		}
	}

	return services, nil
}
