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

const slugUS5 = "datadog-us5"
const nameUS5 = "Datadog US5"

const (
	incidentsUrlUS5           = "https://status.us5.datadoghq.com/api/v2/incidents.json"
	componentsUrlUS5          = "https://status.us5.datadoghq.com/api/v2/components.json"
	scheduleMaintenanceUrlUS5 = "https://status.us5.datadoghq.com/api/v2/scheduled-maintenances.json"
)

type datadogUS5 struct {
	RestyClient *resty.Client
	ServiceID   uint
}

func (d datadogUS5) ID() uint {
	return d.ServiceID
}

func (d datadogUS5) Name() string {
	return nameUS5
}

var _ statuspage.StatusPageProvider = datadogUS5{}

// FetchComponents implements statuspage.Statuspage.
func (d datadogUS5) ScrapComponents() (components.AggregateComponents, error) {
	dComponents := atlassian.ComponentsReq{}
	_, err := d.RestyClient.
		R().
		SetResult(&dComponents).
		EnableRetryDefaultConditions().
		Get(componentsUrlUS5)
	if err != nil {
		return components.AggregateComponents{}, err
	}

	componentGroups := atlassian.FetchComponentsHelper(dComponents)
	return componentGroups, nil
}

// Name implements statuspage.Statuspage.
func (d datadogUS5) Slug() services.ServiceSlug {
	return slugUS5
}

// ScrapIncidents implements statuspage.Statuspage.
func (d datadogUS5) ScrapIncidents() ([]incidents.Incident, error) {
	var req atlassian.IncidentReq
	_, err := d.RestyClient.
		R().
		SetResult(&req).
		Get(incidentsUrlUS5)
	if err != nil {
		return nil, err
	}

	incidents := atlassian.FetchIncidentsHelper(req)

	return incidents, nil
}

// ScrapScheduleMaintainance implements statuspage.Statuspage.
func (d datadogUS5) ScrapScheduleMaintainance() ([]incidents.Incident, error) {
	var req atlassian.IncidentReq
	_, err := d.RestyClient.
		R().
		SetResult(&req).
		Get(scheduleMaintenanceUrlUS5)
	if err != nil {
		return nil, err
	}

	incidents := atlassian.FetchIncidentsHelper(req)

	return incidents, nil
}

func (d datadogUS5) NewWithServiceID(id uint) statuspage.StatusPageProvider {
	d.ServiceID = id
	return d
}

func registerUS5() {
	registry.Register(slugUS5, datadogUS5{
		RestyClient: resty.New(),
	})
}
