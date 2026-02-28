package circleci

import (
	"github.com/yash492/statusy/internal/adapter/collector/atlassian"
	"github.com/yash492/statusy/internal/adapter/collector/registry"
	"github.com/yash492/statusy/internal/domain/components"
	"github.com/yash492/statusy/internal/domain/incidents"
	"github.com/yash492/statusy/internal/domain/services"
	"github.com/yash492/statusy/internal/domain/statuspage"
	"resty.dev/v3"
)

const slug = "circleci"
const name = "Circle CI"

type circleCi struct {
	ScheduleMaintenanceUrl string
	IncidentsUrl           string
	ComponentsUrl          string
	RestyClient            *resty.Client
	ServiceID              uint
}

func (c circleCi) ID() uint {
	return c.ServiceID
}

func (c circleCi) Name() string {
	return name
}

var _ statuspage.StatusPageProvider = circleCi{}

// FetchComponents implements statuspage.Statuspage.
func (c circleCi) ScrapComponents() (components.AggregateComponents, error) {
	circleciComponents := atlassian.ComponentsReq{}
	_, err := c.RestyClient.
		R().
		SetResult(&circleciComponents).
		EnableRetryDefaultConditions().
		Get(c.ComponentsUrl)
	if err != nil {
		return components.AggregateComponents{}, err
	}

	componentGroups := atlassian.FetchComponentsHelper(circleciComponents)
	return componentGroups, nil
}

// Name implements statuspage.Statuspage.
func (c circleCi) Slug() statuspage.ServiceSlug {
	return slug
}

// ScrapIncidents implements statuspage.Statuspage.
func (c circleCi) ScrapIncidents() ([]incidents.Incident, error) {
	var req atlassian.IncidentReq
	_, err := c.RestyClient.
		R().
		SetResult(&req).
		Get(c.IncidentsUrl)
	if err != nil {
		return nil, err
	}

	incidents := atlassian.FetchIncidentsHelper(req)

	return incidents, nil
}

// ScrapScheduleMaintainance implements statuspage.Statuspage.
func (c circleCi) ScrapScheduleMaintainance() ([]incidents.Incident, error) {
	var req atlassian.IncidentReq
	_, err := c.RestyClient.
		R().
		SetResult(&req).
		Get(c.IncidentsUrl)
	if err != nil {
		return nil, err
	}

	incidents := atlassian.FetchIncidentsHelper(req)

	return incidents, nil
}

func Register() {
	registry.Register(slug, func(serviceResult services.ServiceResult) statuspage.StatusPageProvider {
		return circleCi{
			ScheduleMaintenanceUrl: serviceResult.ScheduleMaintenancesUrl,
			IncidentsUrl:           serviceResult.IncidentsUrl,
			ComponentsUrl:          serviceResult.ComponentsUrl,
			RestyClient:            resty.New(),
			ServiceID:              serviceResult.ID,
		}
	})
}
