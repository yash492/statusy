package pleo

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

const slug = "pleo"
const name = "Pleo"

const (
	statusPageUrl = "https://status.pleo.io"
	summaryUrl    = "https://status.pleo.io/proxy/status.pleo.io"
	incidentsUrl  = "https://status.pleo.io/proxy/status.pleo.io/incidents"
)

type pleo struct {
	RestyClient *resty.Client
	ServiceID   uint
}

func (p pleo) ID() uint {
	return p.ServiceID
}

func (p pleo) Name() string {
	return name
}

var _ statuspage.StatusPageProvider = pleo{}

func (p pleo) ScrapComponents() (components.AggregateComponents, error) {
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

func (p pleo) Slug() services.ServiceSlug {
	return slug
}

func (p pleo) ScrapIncidents() ([]incidents.Incident, error) {
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

func (p pleo) ScrapScheduledMaintenance() ([]scheduledmaintenance.ScheduledMaintenance, error) {
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

func (p pleo) GetStatuspageUrl() string {
	return statusPageUrl
}

func (p pleo) NewWithServiceID(id uint) statuspage.StatusPageProvider {
	p.ServiceID = id
	return p
}

func Register(client *resty.Client) {
	registry.Register(slug, pleo{
		RestyClient: client,
	})
}
