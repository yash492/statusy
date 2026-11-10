package harness

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

const slug = "harness"
const name = "Harness"

const (
	incidentsUrl            = "https://status.harness.io/api/v2/incidents.json"
	componentsUrl           = "https://status.harness.io/api/v2/components.json"
	scheduledMaintenanceUrl = "https://status.harness.io/api/v2/scheduled-maintenances.json"
	statusPageUrl           = "https://status.harness.io"
)

type provider struct {
	RestyClient *resty.Client
	ServiceID   uint
}

func (p provider) ID() uint {
	return p.ServiceID
}

func (p provider) Name() string {
	return name
}

var _ statuspage.StatusPageProvider = provider{}

func (p provider) ScrapComponents() (components.AggregateComponents, error) {
	gComponents := atlassian.ComponentsReq{}
	_, err := p.RestyClient.
		R().
		SetResult(&gComponents).
		EnableRetryDefaultConditions().
		Get(componentsUrl)
	if err != nil {
		return components.AggregateComponents{}, err
	}

	return atlassian.FetchComponentsHelper(gComponents), nil
}

func (p provider) Slug() services.ServiceSlug {
	return slug
}

func (p provider) ScrapIncidents() ([]incidents.Incident, error) {
	var req atlassian.IncidentReq
	_, err := p.RestyClient.
		R().
		SetResult(&req).
		Get(incidentsUrl)
	if err != nil {
		return nil, err
	}

	return atlassian.FetchIncidentsHelper(req), nil
}

func (p provider) ScrapScheduledMaintenance() ([]scheduledmaintenance.ScheduledMaintenance, error) {
	var req atlassian.ScheduledMaintenanceReq
	_, err := p.RestyClient.
		R().
		SetResult(&req).
		Get(scheduledMaintenanceUrl)
	if err != nil {
		return nil, err
	}

	return atlassian.FetchScheduledMaintenanceHelper(req), nil
}

func (p provider) GetStatuspageUrl() string {
	return statusPageUrl
}

func (p provider) NewWithServiceID(id uint) statuspage.StatusPageProvider {
	p.ServiceID = id
	return p
}

func Register(client *resty.Client) {
	registry.Register(slug, provider{
		RestyClient: client,
	})
}
