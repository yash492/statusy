package solarwindsobservability

import (
	"github.com/yash492/statusy/internal/adapter/collector/atlassian"
	"github.com/yash492/statusy/internal/adapter/collector/registry"
	"github.com/yash492/statusy/internal/domain/components"
	"github.com/yash492/statusy/internal/domain/incidents"
	"github.com/yash492/statusy/internal/domain/services"
	"github.com/yash492/statusy/internal/domain/statuspage"
	"resty.dev/v3"
)

const slug = "solarwinds-observability"
const name = "SolarWinds Observability"

const (
	incidentsUrl            = "https://status.cloud.solarwinds.com/api/v2/incidents.json"
	componentsUrl           = "https://status.cloud.solarwinds.com/api/v2/components.json"
	scheduledMaintenanceUrl = "https://status.cloud.solarwinds.com/api/v2/scheduled-maintenances.json"
)

type solarWinds struct {
	RestyClient *resty.Client
	ServiceID   uint
}

func (s solarWinds) ID() uint {
	return s.ServiceID
}

func (s solarWinds) Name() string {
	return name
}

var _ statuspage.StatusPageProvider = solarWinds{}

// FetchComponents implements statuspage.Statuspage.
func (s solarWinds) ScrapComponents() (components.AggregateComponents, error) {
	swComponents := atlassian.ComponentsReq{}
	_, err := s.RestyClient.
		R().
		SetResult(&swComponents).
		EnableRetryDefaultConditions().
		Get(componentsUrl)
	if err != nil {
		return components.AggregateComponents{}, err
	}

	componentGroups := atlassian.FetchComponentsHelper(swComponents)
	return componentGroups, nil
}

// Name implements statuspage.Statuspage.
func (s solarWinds) Slug() services.ServiceSlug {
	return slug
}

// ScrapIncidents implements statuspage.Statuspage.
func (s solarWinds) ScrapIncidents() ([]incidents.Incident, error) {
	var req atlassian.IncidentReq
	_, err := s.RestyClient.
		R().
		SetResult(&req).
		Get(incidentsUrl)
	if err != nil {
		return nil, err
	}

	incidents := atlassian.FetchIncidentsHelper(req)

	return incidents, nil
}

// ScrapscheduledMaintenance implements statuspage.Statuspage.
func (s solarWinds) ScrapscheduledMaintenance() ([]incidents.Incident, error) {
	var req atlassian.IncidentReq
	_, err := s.RestyClient.
		R().
		SetResult(&req).
		Get(scheduledMaintenanceUrl)
	if err != nil {
		return nil, err
	}

	incidents := atlassian.FetchIncidentsHelper(req)

	return incidents, nil
}

func (s solarWinds) NewWithServiceID(id uint) statuspage.StatusPageProvider {
	s.ServiceID = id
	return s
}

func Register(client *resty.Client) {
	registry.Register(slug, solarWinds{
		RestyClient: client,
	})
}
