package hashicorp

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

const slug = "hashicorp"
const name = "HashiCorp"

const (
	statusPageUrl = "https://status.hashicorp.com"
	summaryUrl    = "https://status.hashicorp.com/proxy/status.hashicorp.com"
	incidentsUrl  = "https://status.hashicorp.com/proxy/status.hashicorp.com/incidents"
)

type hashicorp struct {
	RestyClient *resty.Client
	ServiceID   uint
}

func (h hashicorp) ID() uint {
	return h.ServiceID
}

func (h hashicorp) Name() string {
	return name
}

var _ statuspage.StatusPageProvider = hashicorp{}

func (h hashicorp) ScrapComponents() (components.AggregateComponents, error) {
	var statusReq incidentio.StatusReq
	_, err := h.RestyClient.
		R().
		SetResult(&statusReq).
		EnableRetryDefaultConditions().
		Get(summaryUrl)
	if err != nil {
		return components.AggregateComponents{}, err
	}

	return incidentio.FetchComponentsHelper(statusReq), nil
}

func (h hashicorp) Slug() services.ServiceSlug {
	return slug
}

func (h hashicorp) ScrapIncidents() ([]incidents.Incident, error) {
	var req incidentio.IncidentsReq
	_, err := h.RestyClient.
		R().
		SetResult(&req).
		Get(incidentsUrl)
	if err != nil {
		return nil, err
	}

	return h.FetchIncidentsHelper(req)
}

func (h hashicorp) FetchIncidentsHelper(req incidentio.IncidentsReq) ([]incidents.Incident, error) {
	return incidentio.FetchIncidentsHelper(req, statusPageUrl), nil
}

func (h hashicorp) ScrapScheduledMaintenance() ([]scheduledmaintenance.ScheduledMaintenance, error) {
	var req incidentio.IncidentsReq
	_, err := h.RestyClient.
		R().
		SetResult(&req).
		Get(incidentsUrl)
	if err != nil {
		return nil, err
	}

	return h.FetchScheduledMaintenancesHelper(req)
}

func (h hashicorp) FetchScheduledMaintenancesHelper(req incidentio.IncidentsReq) ([]scheduledmaintenance.ScheduledMaintenance, error) {
	return incidentio.FetchScheduledMaintenancesHelper(req, statusPageUrl), nil
}

func (h hashicorp) GetStatuspageUrl() string {
	return statusPageUrl
}

func (h hashicorp) NewWithServiceID(id uint) statuspage.StatusPageProvider {
	h.ServiceID = id
	return h
}

func Register(client *resty.Client) {
	registry.Register(slug, hashicorp{
		RestyClient: client,
	})
}
