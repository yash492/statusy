package cursor

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

const slug = "cursor"
const name = "Cursor"

const (
	incidentsUrl            = "https://status.cursor.com/api/v2/incidents.json"
	componentsUrl           = "https://status.cursor.com/api/v2/components.json"
	scheduledMaintenanceUrl = "https://status.cursor.com/api/v2/scheduled-maintenances.json"
)

type cursor struct {
	RestyClient *resty.Client
	ServiceID   uint
}

func (c cursor) ID() uint {
	return c.ServiceID
}

func (c cursor) Name() string {
	return name
}

var _ statuspage.StatusPageProvider = cursor{}

func (c cursor) ScrapComponents() (components.AggregateComponents, error) {
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

func (c cursor) Slug() services.ServiceSlug {
	return slug
}

func (c cursor) ScrapIncidents() ([]incidents.Incident, error) {
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

func (c cursor) ScrapScheduledMaintenance() ([]scheduledmaintenance.ScheduledMaintenance, error) {
	var req atlassian.ScheduledMaintenanceReq
	_, err := c.RestyClient.
		R().
		SetResult(&req).
		Get(scheduledMaintenanceUrl)
	if err != nil {
		return nil, err
	}

	scheduledMaintenances := atlassian.FetchScheduledMaintenanceHelper(req)

	return scheduledMaintenances, nil
}

func (c cursor) NewWithServiceID(id uint) statuspage.StatusPageProvider {
	c.ServiceID = id
	return c
}

func Register(client *resty.Client) {
	registry.Register(slug, cursor{
		RestyClient: client,
	})
}
