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

const slugAP2 = "datadog-ap2"
const nameAP2 = "Datadog AP2"

const (
	incidentsUrlAP2            = "https://status.ap2.datadoghq.com/api/v2/incidents.json"
	componentsUrlAP2           = "https://status.ap2.datadoghq.com/api/v2/components.json"
	scheduledMaintenanceUrlAP2 = "https://status.ap2.datadoghq.com/api/v2/scheduled-maintenances.json"
)

type datadogAP2 struct {
	RestyClient *resty.Client
	ServiceID   uint
}

func (d datadogAP2) ID() uint {
	return d.ServiceID
}

func (d datadogAP2) Name() string {
	return nameAP2
}

var _ statuspage.StatusPageProvider = datadogAP2{}

// FetchComponents implements statuspage.Statuspage.
func (d datadogAP2) ScrapComponents() (components.AggregateComponents, error) {
	dComponents := atlassian.ComponentsReq{}
	_, err := d.RestyClient.
		R().
		SetResult(&dComponents).
		EnableRetryDefaultConditions().
		Get(componentsUrlAP2)
	if err != nil {
		return components.AggregateComponents{}, err
	}

	componentGroups := atlassian.FetchComponentsHelper(dComponents)
	return componentGroups, nil
}

// Name implements statuspage.Statuspage.
func (d datadogAP2) Slug() services.ServiceSlug {
	return slugAP2
}

// ScrapIncidents implements statuspage.Statuspage.
func (d datadogAP2) ScrapIncidents() ([]incidents.Incident, error) {
	var req atlassian.IncidentReq
	_, err := d.RestyClient.
		R().
		SetResult(&req).
		Get(incidentsUrlAP2)
	if err != nil {
		return nil, err
	}

	incidents := atlassian.FetchIncidentsHelper(req)

	return incidents, nil
}

// ScrapscheduledMaintenance implements statuspage.Statuspage.
func (d datadogAP2) ScrapscheduledMaintenance() ([]scheduledmaintenance.ScheduledMaintenance, error) {
	var req atlassian.ScheduledMaintenanceReq
	_, err := d.RestyClient.
		R().
		SetResult(&req).
		Get(scheduledMaintenanceUrlAP2)
	if err != nil {
		return nil, err
	}

	scheduledMaintenances := atlassian.FetchScheduledMaintenanceHelper(req)

	return scheduledMaintenances, nil
}

func (d datadogAP2) NewWithServiceID(id uint) statuspage.StatusPageProvider {
	d.ServiceID = id
	return d
}

func registerAP2(client *resty.Client) {
	registry.Register(slugAP2, datadogAP2{
		RestyClient: client,
	})
}
