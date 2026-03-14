package discord

import (
	"github.com/yash492/statusy/internal/adapter/collector/atlassian"
	"github.com/yash492/statusy/internal/adapter/collector/registry"
	"github.com/yash492/statusy/internal/domain/components"
	"github.com/yash492/statusy/internal/domain/incidents"
	"github.com/yash492/statusy/internal/domain/services"
	"github.com/yash492/statusy/internal/domain/statuspage"
	"resty.dev/v3"
)

const slug = "discord"
const name = "Discord"

const (
	incidentsUrl           = "https://discordstatus.com/api/v2/incidents.json"
	componentsUrl          = "https://discordstatus.com/api/v2/components.json"
	scheduleMaintenanceUrl = "https://discordstatus.com/api/v2/scheduled-maintenances.json"
)

type discord struct {
	RestyClient *resty.Client
	ServiceID   uint
}

func (d discord) ID() uint {
	return d.ServiceID
}

func (d discord) Name() string {
	return name
}

var _ statuspage.StatusPageProvider = discord{}

// FetchComponents implements statuspage.Statuspage.
func (d discord) ScrapComponents() (components.AggregateComponents, error) {
	dComponents := atlassian.ComponentsReq{}
	_, err := d.RestyClient.
		R().
		SetResult(&dComponents).
		EnableRetryDefaultConditions().
		Get(componentsUrl)
	if err != nil {
		return components.AggregateComponents{}, err
	}

	componentGroups := atlassian.FetchComponentsHelper(dComponents)
	return componentGroups, nil
}

// Name implements statuspage.Statuspage.
func (d discord) Slug() services.ServiceSlug {
	return slug
}

// ScrapIncidents implements statuspage.Statuspage.
func (d discord) ScrapIncidents() ([]incidents.Incident, error) {
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

// ScrapScheduleMaintainance implements statuspage.Statuspage.
func (d discord) ScrapScheduleMaintainance() ([]incidents.Incident, error) {
	var req atlassian.IncidentReq
	_, err := d.RestyClient.
		R().
		SetResult(&req).
		Get(scheduleMaintenanceUrl)
	if err != nil {
		return nil, err
	}

	incidents := atlassian.FetchIncidentsHelper(req)

	return incidents, nil
}

func (d discord) NewWithServiceID(id uint) statuspage.StatusPageProvider {
	d.ServiceID = id
	return d
}

func Register() {
	registry.Register(slug, discord{
		RestyClient: resty.New(),
	})
}
