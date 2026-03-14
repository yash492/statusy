package datadog

import (
	"github.com/yash492/statusy/internal/adapter/collector/atlassian"
	"github.com/yash492/statusy/internal/adapter/collector/registry"
	"github.com/yash492/statusy/internal/domain/components"
	"github.com/yash492/statusy/internal/domain/incidents"
	"github.com/yash492/statusy/internal/domain/services"
	"github.com/yash492/statusy/internal/domain/statuspage"
	"resty.dev/v3"
)

const slugEU = "datadog-eu"
const nameEU = "Datadog EU"

const (
	incidentsUrlEU           = "https://status.datadoghq.eu/api/v2/incidents.json"
	componentsUrlEU          = "https://status.datadoghq.eu/api/v2/components.json"
	scheduleMaintenanceUrlEU = "https://status.datadoghq.eu/api/v2/scheduled-maintenances.json"
)

type datadogEU struct {
	RestyClient *resty.Client
	ServiceID   uint
}

func (d datadogEU) ID() uint {
	return d.ServiceID
}

func (d datadogEU) Name() string {
	return nameEU
}

var _ statuspage.StatusPageProvider = datadogEU{}

// FetchComponents implements statuspage.Statuspage.
func (d datadogEU) ScrapComponents() (components.AggregateComponents, error) {
	dComponents := atlassian.ComponentsReq{}
	_, err := d.RestyClient.
		R().
		SetResult(&dComponents).
		EnableRetryDefaultConditions().
		Get(componentsUrlEU)
	if err != nil {
		return components.AggregateComponents{}, err
	}

	componentGroups := atlassian.FetchComponentsHelper(dComponents)
	return componentGroups, nil
}

// Name implements statuspage.Statuspage.
func (d datadogEU) Slug() services.ServiceSlug {
	return slugEU
}

// ScrapIncidents implements statuspage.Statuspage.
func (d datadogEU) ScrapIncidents() ([]incidents.Incident, error) {
	var req atlassian.IncidentReq
	_, err := d.RestyClient.
		R().
		SetResult(&req).
		Get(incidentsUrlEU)
	if err != nil {
		return nil, err
	}

	incidents := atlassian.FetchIncidentsHelper(req)

	return incidents, nil
}

// ScrapScheduleMaintainance implements statuspage.Statuspage.
func (d datadogEU) ScrapScheduleMaintainance() ([]incidents.Incident, error) {
	var req atlassian.IncidentReq
	_, err := d.RestyClient.
		R().
		SetResult(&req).
		Get(scheduleMaintenanceUrlEU)
	if err != nil {
		return nil, err
	}

	incidents := atlassian.FetchIncidentsHelper(req)

	return incidents, nil
}

func (d datadogEU) NewWithServiceID(id uint) statuspage.StatusPageProvider {
	d.ServiceID = id
	return d
}

func RegisterEU() {
	registry.Register(slugEU, datadogEU{
		RestyClient: resty.New(),
	})
}
