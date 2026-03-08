package plivo

import (
	"github.com/yash492/statusy/internal/adapter/collector/atlassian"
	"github.com/yash492/statusy/internal/adapter/collector/registry"
	"github.com/yash492/statusy/internal/domain/components"
	"github.com/yash492/statusy/internal/domain/incidents"
	"github.com/yash492/statusy/internal/domain/services"
	"github.com/yash492/statusy/internal/domain/statuspage"
	"resty.dev/v3"
)

const slug = "plivo"
const name = "Plivo"

const (
	incidentsUrl           = "https://status.plivo.com/api/v2/incidents.json"
	componentsUrl          = "https://status.plivo.com/api/v2/components.json"
	scheduleMaintenanceUrl = "https://status.plivo.com/api/v2/scheduled-maintenances.json"
)

type plivo struct {
	RestyClient *resty.Client
	ServiceID   uint
}

func (p plivo) ID() uint {
	return p.ServiceID
}

func (p plivo) Name() string {
	return name
}

var _ statuspage.StatusPageProvider = plivo{}

func (p plivo) ScrapComponents() (components.AggregateComponents, error) {
	plivoComponents := atlassian.ComponentsReq{}
	_, err := p.RestyClient.
		R().
		SetResult(&plivoComponents).
		EnableRetryDefaultConditions().
		Get(componentsUrl)
	if err != nil {
		return components.AggregateComponents{}, err
	}
	return atlassian.FetchComponentsHelper(plivoComponents), nil
}

func (p plivo) Slug() services.ServiceSlug {
	return slug
}

func (p plivo) ScrapIncidents() ([]incidents.Incident, error) {
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

func (p plivo) ScrapScheduleMaintainance() ([]incidents.Incident, error) {
	var req atlassian.IncidentReq
	_, err := p.RestyClient.
		R().
		SetResult(&req).
		Get(scheduleMaintenanceUrl)
	if err != nil {
		return nil, err
	}
	return atlassian.FetchIncidentsHelper(req), nil
}

func (p plivo) NewWithServiceID(id uint) statuspage.StatusPageProvider {
	p.ServiceID = id
	return p
}

func Register() {
	registry.Register(slug, plivo{
		RestyClient: resty.New(),
	})
}
