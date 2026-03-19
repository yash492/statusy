package newrelic

import (
	"github.com/yash492/statusy/internal/adapter/collector/atlassian"
	"github.com/yash492/statusy/internal/adapter/collector/registry"
	"github.com/yash492/statusy/internal/domain/components"
	"github.com/yash492/statusy/internal/domain/incidents"
	"github.com/yash492/statusy/internal/domain/services"
	"github.com/yash492/statusy/internal/domain/statuspage"
	"resty.dev/v3"
)

const slug = "newrelic"
const name = "New Relic"

const (
	incidentsUrl           = "https://status.newrelic.com/api/v2/incidents.json"
	componentsUrl          = "https://status.newrelic.com/api/v2/components.json"
	scheduleMaintenanceUrl = "https://status.newrelic.com/api/v2/scheduled-maintenances.json"
)

type newrelic struct {
	RestyClient *resty.Client
	ServiceID   uint
}

func (n newrelic) ID() uint {
	return n.ServiceID
}

func (n newrelic) Name() string {
	return name
}

var _ statuspage.StatusPageProvider = newrelic{}

// FetchComponents implements statuspage.Statuspage.
func (n newrelic) ScrapComponents() (components.AggregateComponents, error) {
	nComponents := atlassian.ComponentsReq{}
	_, err := n.RestyClient.
		R().
		SetResult(&nComponents).
		EnableRetryDefaultConditions().
		Get(componentsUrl)
	if err != nil {
		return components.AggregateComponents{}, err
	}

	componentGroups := atlassian.FetchComponentsHelper(nComponents)
	return componentGroups, nil
}

// Name implements statuspage.Statuspage.
func (n newrelic) Slug() services.ServiceSlug {
	return slug
}

// ScrapIncidents implements statuspage.Statuspage.
func (n newrelic) ScrapIncidents() ([]incidents.Incident, error) {
	var req atlassian.IncidentReq
	_, err := n.RestyClient.
		R().
		SetResult(&req).
		Get(incidentsUrl)
	if err != nil {
		return nil, err
	}

	incidents := atlassian.FetchIncidentsHelper(req)

	return incidents, nil
}

// ScrapScheduleMaintainance implements statuspage.Statuspage.
func (n newrelic) ScrapScheduleMaintainance() ([]incidents.Incident, error) {
	var req atlassian.IncidentReq
	_, err := n.RestyClient.
		R().
		SetResult(&req).
		Get(scheduleMaintenanceUrl)
	if err != nil {
		return nil, err
	}

	incidents := atlassian.FetchIncidentsHelper(req)

	return incidents, nil
}

func (n newrelic) NewWithServiceID(id uint) statuspage.StatusPageProvider {
	n.ServiceID = id
	return n
}

func Register(client *resty.Client) {
	registry.Register(slug, newrelic{
		RestyClient: client,
	})
}
