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

const slugGov = "datadog-gov"
const nameGov = "Datadog GovCloud"

const (
	incidentsUrlGov            = "https://status.ddog-gov.com/api/v2/incidents.json"
	componentsUrlGov           = "https://status.ddog-gov.com/api/v2/components.json"
	scheduledMaintenanceUrlGov = "https://status.ddog-gov.com/api/v2/scheduled-maintenances.json"
)

type datadogGov struct {
	RestyClient *resty.Client
	ServiceID   uint
}

func (d datadogGov) ID() uint {
	return d.ServiceID
}

func (d datadogGov) Name() string {
	return nameGov
}

var _ statuspage.StatusPageProvider = datadogGov{}

// FetchComponents implements statuspage.Statuspage.
func (d datadogGov) ScrapComponents() (components.AggregateComponents, error) {
	dComponents := atlassian.ComponentsReq{}
	_, err := d.RestyClient.
		R().
		SetResult(&dComponents).
		EnableRetryDefaultConditions().
		Get(componentsUrlGov)
	if err != nil {
		return components.AggregateComponents{}, err
	}

	componentGroups := atlassian.FetchComponentsHelper(dComponents)
	return componentGroups, nil
}

// Name implements statuspage.Statuspage.
func (d datadogGov) Slug() services.ServiceSlug {
	return slugGov
}

// ScrapIncidents implements statuspage.Statuspage.
func (d datadogGov) ScrapIncidents() ([]incidents.Incident, error) {
	var req atlassian.IncidentReq
	_, err := d.RestyClient.
		R().
		SetResult(&req).
		Get(incidentsUrlGov)
	if err != nil {
		return nil, err
	}

	incidents := atlassian.FetchIncidentsHelper(req)

	return incidents, nil
}

// ScrapscheduledMaintenance implements statuspage.Statuspage.
func (d datadogGov) ScrapscheduledMaintenance() ([]scheduledmaintenance.ScheduledMaintenance, error) {
	var req atlassian.ScheduledMaintenanceReq
	_, err := d.RestyClient.
		R().
		SetResult(&req).
		Get(scheduledMaintenanceUrlGov)
	if err != nil {
		return nil, err
	}

	scheduledMaintenances := atlassian.FetchScheduledMaintenanceHelper(req)

	return scheduledMaintenances, nil
}

func (d datadogGov) NewWithServiceID(id uint) statuspage.StatusPageProvider {
	d.ServiceID = id
	return d
}

func registerGov(client *resty.Client) {
	registry.Register(slugGov, datadogGov{
		RestyClient: client,
	})
}
