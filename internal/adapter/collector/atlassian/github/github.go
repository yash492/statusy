package github

import (
	"github.com/yash492/statusy/internal/adapter/collector/atlassian"
	"github.com/yash492/statusy/internal/adapter/collector/registry"
	"github.com/yash492/statusy/internal/domain/components"
	"github.com/yash492/statusy/internal/domain/incidents"
	"github.com/yash492/statusy/internal/domain/services"
	"github.com/yash492/statusy/internal/domain/statuspage"
	"resty.dev/v3"
)

const slug = "github"
const name = "GitHub"

const (
	incidentsUrl            = "https://www.githubstatus.com/api/v2/incidents.json"
	componentsUrl           = "https://www.githubstatus.com/api/v2/components.json"
	scheduledMaintenanceUrl = "https://www.githubstatus.com/api/v2/scheduled-maintenances.json"
)

type github struct {
	RestyClient *resty.Client
	ServiceID   uint
}

func (g github) ID() uint {
	return g.ServiceID
}

func (g github) Name() string {
	return name
}

var _ statuspage.StatusPageProvider = github{}

// FetchComponents implements statuspage.Statuspage.
func (g github) ScrapComponents() (components.AggregateComponents, error) {
	gComponents := atlassian.ComponentsReq{}
	_, err := g.RestyClient.
		R().
		SetResult(&gComponents).
		EnableRetryDefaultConditions().
		Get(componentsUrl)
	if err != nil {
		return components.AggregateComponents{}, err
	}

	componentGroups := atlassian.FetchComponentsHelper(gComponents)
	return componentGroups, nil
}

// Name implements statuspage.Statuspage.
func (g github) Slug() services.ServiceSlug {
	return slug
}

// ScrapIncidents implements statuspage.Statuspage.
func (g github) ScrapIncidents() ([]incidents.Incident, error) {
	var req atlassian.IncidentReq
	_, err := g.RestyClient.
		R().
		SetResult(&req).
		Get(incidentsUrl)
	if err != nil {
		return nil, err
	}

	incidents := atlassian.FetchIncidentsHelper(req)

	return incidents, nil
}

// ScrapscheduledMaintenance implements statuspage.Statuspage.
func (g github) ScrapscheduledMaintenance() ([]incidents.Incident, error) {
	var req atlassian.IncidentReq
	_, err := g.RestyClient.
		R().
		SetResult(&req).
		Get(scheduledMaintenanceUrl)
	if err != nil {
		return nil, err
	}

	incidents := atlassian.FetchIncidentsHelper(req)

	return incidents, nil
}

func (g github) NewWithServiceID(id uint) statuspage.StatusPageProvider {
	g.ServiceID = id
	return g
}

func Register(client *resty.Client) {
	registry.Register(slug, github{
		RestyClient: client,
	})
}
