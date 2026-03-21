package zoom

import (
	"github.com/yash492/statusy/internal/adapter/collector/atlassian"
	"github.com/yash492/statusy/internal/adapter/collector/registry"
	"github.com/yash492/statusy/internal/domain/components"
	"github.com/yash492/statusy/internal/domain/incidents"
	"github.com/yash492/statusy/internal/domain/services"
	"github.com/yash492/statusy/internal/domain/statuspage"
	"resty.dev/v3"
)

const slug = "zoom"
const name = "Zoom"

const (
	incidentsUrl            = "https://www.zoomstatus.com/api/v2/incidents.json"
	componentsUrl           = "https://www.zoomstatus.com/api/v2/components.json"
	scheduledMaintenanceUrl = "https://www.zoomstatus.com/api/v2/scheduled-maintenances.json"
)

type zoomProvider struct {
	RestyClient *resty.Client
	ServiceID   uint
}

func (z zoomProvider) ID() uint {
	return z.ServiceID
}

func (z zoomProvider) Name() string {
	return name
}

var _ statuspage.StatusPageProvider = zoomProvider{}

// ScrapComponents implements statuspage.Statuspage.
func (z zoomProvider) ScrapComponents() (components.AggregateComponents, error) {
	zmComponents := atlassian.ComponentsReq{}
	_, err := z.RestyClient.
		R().
		SetResult(&zmComponents).
		EnableRetryDefaultConditions().
		Get(componentsUrl)
	if err != nil {
		return components.AggregateComponents{}, err
	}

	componentGroups := atlassian.FetchComponentsHelper(zmComponents)
	return componentGroups, nil
}

// Slug implements statuspage.Statuspage.
func (z zoomProvider) Slug() services.ServiceSlug {
	return slug
}

// ScrapIncidents implements statuspage.Statuspage.
func (z zoomProvider) ScrapIncidents() ([]incidents.Incident, error) {
	var req atlassian.IncidentReq
	_, err := z.RestyClient.
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
func (z zoomProvider) ScrapscheduledMaintenance() ([]incidents.Incident, error) {
	var req atlassian.IncidentReq
	_, err := z.RestyClient.
		R().
		SetResult(&req).
		Get(scheduledMaintenanceUrl)
	if err != nil {
		return nil, err
	}

	incidents := atlassian.FetchIncidentsHelper(req)

	return incidents, nil
}

func (z zoomProvider) NewWithServiceID(id uint) statuspage.StatusPageProvider {
	z.ServiceID = id
	return z
}

func Register(client *resty.Client) {
	registry.Register(slug, zoomProvider{
		RestyClient: client,
	})
}
