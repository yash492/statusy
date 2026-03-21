package digitalocean

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

const slug = "digitalocean"
const name = "Digital Ocean"

const (
	incidentsUrl            = "https://status.digitalocean.com/api/v2/incidents.json"
	componentsUrl           = "https://status.digitalocean.com/api/v2/components.json"
	scheduledMaintenanceUrl = "https://status.digitalocean.com/api/v2/scheduled-maintenances.json"
)

type digitalOcean struct {
	RestyClient *resty.Client
	ServiceID   uint
}

func (d digitalOcean) ID() uint {
	return d.ServiceID
}

func (d digitalOcean) Name() string {
	return name
}

var _ statuspage.StatusPageProvider = digitalOcean{}

// FetchComponents implements statuspage.Statuspage.
func (d digitalOcean) ScrapComponents() (components.AggregateComponents, error) {
	doComponents := atlassian.ComponentsReq{}
	_, err := d.RestyClient.
		R().
		SetResult(&doComponents).
		EnableRetryDefaultConditions().
		Get(componentsUrl)
	if err != nil {
		return components.AggregateComponents{}, err
	}

	componentGroups := atlassian.FetchComponentsHelper(doComponents)
	return componentGroups, nil
}

// Name implements statuspage.Statuspage.
func (d digitalOcean) Slug() services.ServiceSlug {
	return slug
}

// ScrapIncidents implements statuspage.Statuspage.
func (d digitalOcean) ScrapIncidents() ([]incidents.Incident, error) {
	var req atlassian.IncidentReq
	_, err := d.RestyClient.
		R().
		SetResult(&req).
		Get(incidentsUrl)
	if err != nil {
		return nil, err
	}

	incidents := atlassian.FetchIncidentsHelper(req)

	return incidents, nil
}

// ScrapScheduledMaintenance implements statuspage.Statuspage.
func (d digitalOcean) ScrapScheduledMaintenance() ([]scheduledmaintenance.ScheduledMaintenance, error) {
	var req atlassian.ScheduledMaintenanceReq
	_, err := d.RestyClient.
		R().
		SetResult(&req).
		Get(scheduledMaintenanceUrl)
	if err != nil {
		return nil, err
	}

	scheduledMaintenances := atlassian.FetchScheduledMaintenanceHelper(req)

	return scheduledMaintenances, nil
}

func (d digitalOcean) NewWithServiceID(id uint) statuspage.StatusPageProvider {
	d.ServiceID = id
	return d
}

func Register(client *resty.Client) {
	registry.Register(slug, digitalOcean{
		RestyClient: client,
	})
}
