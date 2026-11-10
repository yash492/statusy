package fireworksai

import (
	"github.com/yash492/statusy/internal/adapter/collector/incidentio"
	"github.com/yash492/statusy/internal/adapter/collector/registry"
	"github.com/yash492/statusy/internal/domain/components"
	"github.com/yash492/statusy/internal/domain/incidents"
	"github.com/yash492/statusy/internal/domain/scheduledmaintenance"
	"github.com/yash492/statusy/internal/domain/services"
	"github.com/yash492/statusy/internal/domain/statuspage"
	"resty.dev/v3"
)

const slug = "fireworksai"
const name = "Fireworks AI"

const (
	statusPageUrl = "https://status.fireworks.ai"
	summaryUrl    = "https://status.fireworks.ai/proxy/status.fireworks.ai"
	incidentsUrl  = "https://status.fireworks.ai/proxy/status.fireworks.ai/incidents"
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
	var statusReq incidentio.StatusReq
	_, err := p.RestyClient.
		R().
		SetResult(&statusReq).
		EnableRetryDefaultConditions().
		Get(summaryUrl)
	if err != nil {
		return components.AggregateComponents{}, err
	}

	return incidentio.FetchComponentsHelper(statusReq), nil
}

func (p provider) Slug() services.ServiceSlug {
	return slug
}

func (p provider) ScrapIncidents() ([]incidents.Incident, error) {
	var req incidentio.IncidentsReq
	_, err := p.RestyClient.
		R().
		SetResult(&req).
		Get(incidentsUrl)
	if err != nil {
		return nil, err
	}

	return incidentio.FetchIncidentsHelper(req, statusPageUrl), nil
}

func (p provider) ScrapScheduledMaintenance() ([]scheduledmaintenance.ScheduledMaintenance, error) {
	var req incidentio.IncidentsReq
	_, err := p.RestyClient.
		R().
		SetResult(&req).
		Get(incidentsUrl)
	if err != nil {
		return nil, err
	}

	return incidentio.FetchScheduledMaintenancesHelper(req, statusPageUrl), nil
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
