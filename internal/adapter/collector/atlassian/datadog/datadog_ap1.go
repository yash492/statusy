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

const slugAP1 = "datadog-ap1"
const nameAP1 = "Datadog AP1"

const (
	incidentsUrlAP1            = "https://status.ap1.datadoghq.com/api/v2/incidents.json"
	componentsUrlAP1           = "https://status.ap1.datadoghq.com/api/v2/components.json"
	scheduledMaintenanceUrlAP1 = "https://status.ap1.datadoghq.com/api/v2/scheduled-maintenances.json"
)

type datadogAP1 struct {
	RestyClient *resty.Client
	ServiceID   uint
}

func (d datadogAP1) ID() uint {
	return d.ServiceID
}

func (d datadogAP1) Name() string {
	return nameAP1
}

var _ statuspage.StatusPageProvider = datadogAP1{}

// FetchComponents implements statuspage.Statuspage.
func (d datadogAP1) ScrapComponents() (components.AggregateComponents, error) {
	dComponents := atlassian.ComponentsReq{}
	_, err := d.RestyClient.
		R().
		SetResult(&dComponents).
		EnableRetryDefaultConditions().
		Get(componentsUrlAP1)
	if err != nil {
		return components.AggregateComponents{}, err
	}

	componentGroups := atlassian.FetchComponentsHelper(dComponents)
	return componentGroups, nil
}

// Name implements statuspage.Statuspage.
func (d datadogAP1) Slug() services.ServiceSlug {
	return slugAP1
}

// ScrapIncidents implements statuspage.Statuspage.
func (d datadogAP1) ScrapIncidents() ([]incidents.Incident, error) {
	var req atlassian.IncidentReq
	_, err := d.RestyClient.
		R().
		SetResult(&req).
		Get(incidentsUrlAP1)
	if err != nil {
		return nil, err
	}

	incidents := atlassian.FetchIncidentsHelper(req)

	return incidents, nil
}

// ScrapscheduledMaintenance implements statuspage.Statuspage.
func (d datadogAP1) ScrapscheduledMaintenance() ([]incidents.Incident, error) {
	var req atlassian.IncidentReq
	_, err := d.RestyClient.
		R().
		SetResult(&req).
		Get(scheduledMaintenanceUrlAP1)
	if err != nil {
		return nil, err
	}

	incidents := atlassian.FetchIncidentsHelper(req)

	return incidents, nil
}

func (d datadogAP1) NewWithServiceID(id uint) statuspage.StatusPageProvider {
	d.ServiceID = id
	return d
}

func registerAP1(client *resty.Client) {
	registry.Register(slugAP1, datadogAP1{
		RestyClient: client,
	})
}
