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

const slugUS3 = "datadog-us3"
const nameUS3 = "Datadog US3"

const (
	incidentsUrlUS3           = "https://status.us3.datadoghq.com/api/v2/incidents.json"
	componentsUrlUS3          = "https://status.us3.datadoghq.com/api/v2/components.json"
	scheduleMaintenanceUrlUS3 = "https://status.us3.datadoghq.com/api/v2/scheduled-maintenances.json"
)

type datadogUS3 struct {
	RestyClient *resty.Client
	ServiceID   uint
}

func (d datadogUS3) ID() uint {
	return d.ServiceID
}

func (d datadogUS3) Name() string {
	return nameUS3
}

var _ statuspage.StatusPageProvider = datadogUS3{}

// FetchComponents implements statuspage.Statuspage.
func (d datadogUS3) ScrapComponents() (components.AggregateComponents, error) {
	dComponents := atlassian.ComponentsReq{}
	_, err := d.RestyClient.
		R().
		SetResult(&dComponents).
		EnableRetryDefaultConditions().
		Get(componentsUrlUS3)
	if err != nil {
		return components.AggregateComponents{}, err
	}

	componentGroups := atlassian.FetchComponentsHelper(dComponents)
	return componentGroups, nil
}

// Name implements statuspage.Statuspage.
func (d datadogUS3) Slug() services.ServiceSlug {
	return slugUS3
}

// ScrapIncidents implements statuspage.Statuspage.
func (d datadogUS3) ScrapIncidents() ([]incidents.Incident, error) {
	var req atlassian.IncidentReq
	_, err := d.RestyClient.
		R().
		SetResult(&req).
		Get(incidentsUrlUS3)
	if err != nil {
		return nil, err
	}

	incidents := atlassian.FetchIncidentsHelper(req)

	return incidents, nil
}

// ScrapScheduleMaintainance implements statuspage.Statuspage.
func (d datadogUS3) ScrapScheduleMaintainance() ([]incidents.Incident, error) {
	var req atlassian.IncidentReq
	_, err := d.RestyClient.
		R().
		SetResult(&req).
		Get(scheduleMaintenanceUrlUS3)
	if err != nil {
		return nil, err
	}

	incidents := atlassian.FetchIncidentsHelper(req)

	return incidents, nil
}

func (d datadogUS3) NewWithServiceID(id uint) statuspage.StatusPageProvider {
	d.ServiceID = id
	return d
}

func RegisterUS3() {
	registry.Register(slugUS3, datadogUS3{
		RestyClient: resty.New(),
	})
}
