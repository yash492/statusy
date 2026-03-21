package datadog

import (
	"github.com/yash492/statusy/internal/adapter/collector/atlassian"
	"github.com/yash492/statusy/internal/adapter/collector/registry"
	"github.com/yash492/statusy/internal/domain/components"
	"github.com/yash492/statusy/internal/domain/incidents"
	"github.com/yash492/statusy/internal/domain/scheduledmaintenance"
	"github.com/yash492/statusy/internal/domain/services"
	"github.com/yash492/statusy/internal/domain/statuspage"
	"resty.dev/v3"
)

const slugEU = "datadog-eu"
const nameEU = "Datadog EU"

const (
	incidentsUrlEU            = "https://status.datadoghq.eu/api/v2/incidents.json"
	componentsUrlEU           = "https://status.datadoghq.eu/api/v2/components.json"
	scheduledMaintenanceUrlEU = "https://status.datadoghq.eu/api/v2/scheduled-maintenances.json"
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

// ScrapScheduledMaintenance implements statuspage.Statuspage.
func (d datadogEU) ScrapScheduledMaintenance() ([]scheduledmaintenance.ScheduledMaintenance, error) {
	var req atlassian.ScheduledMaintenanceReq
	_, err := d.RestyClient.
		R().
		SetResult(&req).
		Get(scheduledMaintenanceUrlEU)
	if err != nil {
		return nil, err
	}

	scheduledMaintenances := atlassian.FetchScheduledMaintenanceHelper(req)

	return scheduledMaintenances, nil
}

func (d datadogEU) NewWithServiceID(id uint) statuspage.StatusPageProvider {
	d.ServiceID = id
	return d
}

func registerEU(client *resty.Client) {
	registry.Register(slugEU, datadogEU{
		RestyClient: client,
	})
}
