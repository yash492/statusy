package atlassianstatuspage

import (
	"github.com/yash492/statusy/internal/domain/components"
	"github.com/yash492/statusy/internal/domain/incidents"
	"github.com/yash492/statusy/internal/domain/statuspage"
	"resty.dev/v3"
)

type CircleCi struct {
	ScheduleMaintenanceUrl string
	IncidentsUrl           string
	ComponentsUrl          string
	RestyClient            *resty.Client
	ServiceID              uint
}

func (c CircleCi) ID() uint {
	return c.ServiceID
}

func (c CircleCi) Name() string {
	return "Circle CI"
}

func NewCircleCIProvider(
	componentsUrl string,
	incidentsUrl string,
	scheduleMaintenanceUrl string,
	serviceID uint,
	restyClient *resty.Client,
) CircleCi {
	return CircleCi{
		IncidentsUrl:           incidentsUrl,
		ComponentsUrl:          componentsUrl,
		ScheduleMaintenanceUrl: scheduleMaintenanceUrl,
		RestyClient:            restyClient,
		ServiceID:              serviceID,
	}
}

var _ statuspage.StatusPageProvider = CircleCi{}

// FetchComponents implements statuspage.Statuspage.
func (c CircleCi) ScrapComponents() (components.AggregateComponents, error) {
	circleciComponents := atlassianComponentsReq{}
	_, err := c.RestyClient.
		R().
		SetResult(&circleciComponents).
		EnableRetryDefaultConditions().
		Get(c.ComponentsUrl)
	if err != nil {
		return components.AggregateComponents{}, err
	}

	componentGroups := fetchComponentsHelper(circleciComponents)
	return componentGroups, nil
}

// Name implements statuspage.Statuspage.
func (c CircleCi) Slug() statuspage.ServiceSlug {
	return circleci
}

// ScrapIncidents implements statuspage.Statuspage.
func (c CircleCi) ScrapIncidents() ([]incidents.Incident, error) {
	var req atlassianIncidentReq
	_, err := c.RestyClient.
		R().
		SetResult(&req).
		Get(c.IncidentsUrl)
	if err != nil {
		return nil, err
	}

	incidents := fetchIncidentsHelper(req, c.Slug().String())

	return incidents, nil
}

// ScrapScheduleMaintainance implements statuspage.Statuspage.
func (c CircleCi) ScrapScheduleMaintainance() ([]incidents.Incident, error) {
	var req atlassianIncidentReq
	_, err := c.RestyClient.
		R().
		SetResult(&req).
		Get(c.IncidentsUrl)
	if err != nil {
		return nil, err
	}

	incidents := fetchIncidentsHelper(req, c.Slug().String())

	return incidents, nil
}
