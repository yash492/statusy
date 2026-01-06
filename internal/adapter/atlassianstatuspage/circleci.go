package atlassianstatuspage

import (
	"github.com/yash492/statusy/internal/domain/statuspage"
	"resty.dev/v3"
)

const circleciSlug = "circleci"

type CircleCi struct {
	ScheduleMaintenanceUrl string
	IncidentsUrl           string
	ComponentsUrl          string
	RestyClient            *resty.Client
}

func New(
	componentsUrl string,
	incidentsUrl string,
	scheduleMaintenanceUrl string,
	restyClient *resty.Client,
) CircleCi {
	return CircleCi{
		IncidentsUrl:           incidentsUrl,
		ComponentsUrl:          componentsUrl,
		ScheduleMaintenanceUrl: scheduleMaintenanceUrl,
		RestyClient:            restyClient,
	}
}

var _ statuspage.StatusPageProvider = CircleCi{}

// FetchComponents implements statuspage.Statuspage.
func (c CircleCi) FetchComponents() (statuspage.AggregateComponents, error) {
	circleciComponents := atlassianComponentsReq{}
	_, err := c.RestyClient.
		R().
		SetResult(&circleciComponents).
		EnableRetryDefaultConditions().
		Get(c.ComponentsUrl)
	if err != nil {
		return statuspage.AggregateComponents{}, err
	}

	componentGroups := fetchComponentsHelper(circleciComponents, c.Slug())
	return componentGroups, nil
}

// Name implements statuspage.Statuspage.
func (c CircleCi) Slug() string {
	return circleciSlug
}

// ScrapIncidents implements statuspage.Statuspage.
func (c CircleCi) ScrapIncidents() ([]statuspage.Incident, error) {
	var req atlassianIncidentReq
	_, err := c.RestyClient.
		R().
		SetResult(&req).
		Get(c.IncidentsUrl)
	if err != nil {
		return nil, err
	}

	incidents := fetchIncidentsHelper(req, c.Slug())

	return incidents, nil
}

// ScrapScheduleMaintainance implements statuspage.Statuspage.
func (c CircleCi) ScrapScheduleMaintainance() ([]statuspage.Incident, error) {
	var req atlassianIncidentReq
	_, err := c.RestyClient.
		R().
		SetResult(&req).
		Get(c.IncidentsUrl)
	if err != nil {
		return nil, err
	}

	incidents := fetchIncidentsHelper(req, c.Slug())

	return incidents, nil
}
