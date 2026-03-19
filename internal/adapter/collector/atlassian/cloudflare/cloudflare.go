package cloudflare

import (
	"github.com/yash492/statusy/internal/adapter/collector/atlassian"
	"github.com/yash492/statusy/internal/adapter/collector/registry"
	"github.com/yash492/statusy/internal/domain/components"
	"github.com/yash492/statusy/internal/domain/incidents"
	"github.com/yash492/statusy/internal/domain/services"
	"github.com/yash492/statusy/internal/domain/statuspage"
	"resty.dev/v3"
)

const slug = "cloudflare"
const name = "Cloudflare"

const (
	incidentsUrl           = "https://www.cloudflarestatus.com/api/v2/incidents.json"
	componentsUrl          = "https://www.cloudflarestatus.com/api/v2/components.json"
	scheduleMaintenanceUrl = "https://www.cloudflarestatus.com/api/v2/scheduled-maintenances.json"
)

type cloudflare struct {
	RestyClient *resty.Client
	ServiceID   uint
}

func (c cloudflare) ID() uint {
	return c.ServiceID
}

func (c cloudflare) Name() string {
	return name
}

var _ statuspage.StatusPageProvider = cloudflare{}

// FetchComponents implements statuspage.Statuspage.
func (c cloudflare) ScrapComponents() (components.AggregateComponents, error) {
	cComponents := atlassian.ComponentsReq{}
	_, err := c.RestyClient.
		R().
		SetResult(&cComponents).
		EnableRetryDefaultConditions().
		Get(componentsUrl)
	if err != nil {
		return components.AggregateComponents{}, err
	}

	componentGroups := atlassian.FetchComponentsHelper(cComponents)
	return componentGroups, nil
}

// Name implements statuspage.Statuspage.
func (c cloudflare) Slug() services.ServiceSlug {
	return slug
}

// ScrapIncidents implements statuspage.Statuspage.
func (c cloudflare) ScrapIncidents() ([]incidents.Incident, error) {
	var req atlassian.IncidentReq
	_, err := c.RestyClient.
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
func (c cloudflare) ScrapScheduleMaintainance() ([]incidents.Incident, error) {
	var req atlassian.IncidentReq
	_, err := c.RestyClient.
		R().
		SetResult(&req).
		Get(scheduleMaintenanceUrl)
	if err != nil {
		return nil, err
	}

	incidents := atlassian.FetchIncidentsHelper(req)

	return incidents, nil
}

func (c cloudflare) NewWithServiceID(id uint) statuspage.StatusPageProvider {
	c.ServiceID = id
	return c
}

func Register(client *resty.Client) {
	registry.Register(slug, cloudflare{
		RestyClient: client,
	})
}
