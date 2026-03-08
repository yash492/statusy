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

const (
	incidentsUrl           = "https://status.circleci.com/api/v2/incidents.json"
	componentsUrl          = "https://status.circleci.com/api/v2/components.json"
	scheduleMaintenanceUrl = "https://status.circleci.com/api/v2/scheduled-maintenances.json"
)

type circleCi struct {
	RestyClient *resty.Client
	ServiceID   uint
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
		Get(componentsUrl)
	if err != nil {
		return components.AggregateComponents{}, err
	}

	componentGroups := atlassian.FetchComponentsHelper(circleciComponents)
	return componentGroups, nil
}

// Name implements statuspage.Statuspage.
func (c circleCi) Slug() services.ServiceSlug {
	return slug
}

// ScrapIncidents implements statuspage.Statuspage.
func (c circleCi) ScrapIncidents() ([]incidents.Incident, error) {
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
func (c circleCi) ScrapScheduleMaintainance() ([]incidents.Incident, error) {
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

func Register() {
	registry.Register(slug, func(serviceID uint) statuspage.StatusPageProvider {
		return circleCi{
			RestyClient: resty.New(),
			ServiceID:   serviceID,
		}
	})
}
