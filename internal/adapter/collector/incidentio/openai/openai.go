package openai

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

const slug = "openai"
const name = "OpenAI"

const (
	statusPageUrl = "https://status.openai.com"
	summaryUrl    = "https://status.openai.com/proxy/status.openai.com"
	incidentsUrl  = "https://status.openai.com/proxy/status.openai.com/incidents"
)

type openai struct {
	RestyClient *resty.Client
	ServiceID   uint
}

func (o openai) ID() uint {
	return o.ServiceID
}

func (o openai) Name() string {
	return name
}

var _ statuspage.StatusPageProvider = openai{}

func (o openai) ScrapComponents() (components.AggregateComponents, error) {
	var statusReq incidentio.StatusReq
	_, err := o.RestyClient.
		R().
		SetResult(&statusReq).
		EnableRetryDefaultConditions().
		Get(summaryUrl)
	if err != nil {
		return components.AggregateComponents{}, err
	}

	return incidentio.FetchComponentsHelper(statusReq), nil
}

func (o openai) Slug() services.ServiceSlug {
	return slug
}

func (o openai) ScrapIncidents() ([]incidents.Incident, error) {
	var req incidentio.IncidentsReq
	_, err := o.RestyClient.
		R().
		SetResult(&req).
		Get(incidentsUrl)
	if err != nil {
		return nil, err
	}

	return incidentio.FetchIncidentsHelper(req, statusPageUrl), nil
}

func (o openai) ScrapScheduledMaintenance() ([]scheduledmaintenance.ScheduledMaintenance, error) {
	// Scheduled maintenance is skipped for now as requested
	return []scheduledmaintenance.ScheduledMaintenance{}, nil
}

func (o openai) GetStatuspageUrl() string {
	return statusPageUrl
}

func (o openai) NewWithServiceID(id uint) statuspage.StatusPageProvider {
	o.ServiceID = id
	return o
}

func Register(client *resty.Client) {
	registry.Register(slug, openai{
		RestyClient: client,
	})
}
