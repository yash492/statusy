package twilio

import (
	"github.com/yash492/statusy/internal/adapter/collector/atlassian"
	"github.com/yash492/statusy/internal/adapter/collector/registry"
	"github.com/yash492/statusy/internal/domain/components"
	"github.com/yash492/statusy/internal/domain/incidents"
	"github.com/yash492/statusy/internal/domain/services"
	"github.com/yash492/statusy/internal/domain/statuspage"
	"resty.dev/v3"
)

const slug = "twilio"
const name = "Twilio"

const (
	incidentsUrl           = "https://status.twilio.com/api/v2/incidents.json"
	componentsUrl          = "https://status.twilio.com/api/v2/components.json"
	scheduleMaintenanceUrl = "https://status.twilio.com/api/v2/scheduled-maintenances.json"
)

type twilio struct {
	RestyClient *resty.Client
	ServiceID   uint
}

func (t twilio) ID() uint {
	return t.ServiceID
}

func (t twilio) Name() string {
	return name
}

var _ statuspage.StatusPageProvider = twilio{}

// FetchComponents implements statuspage.Statuspage.
func (t twilio) ScrapComponents() (components.AggregateComponents, error) {
	tComponents := atlassian.ComponentsReq{}
	_, err := t.RestyClient.
		R().
		SetResult(&tComponents).
		EnableRetryDefaultConditions().
		Get(componentsUrl)
	if err != nil {
		return components.AggregateComponents{}, err
	}

	componentGroups := atlassian.FetchComponentsHelper(tComponents)
	return componentGroups, nil
}

// Name implements statuspage.Statuspage.
func (t twilio) Slug() services.ServiceSlug {
	return slug
}

// ScrapIncidents implements statuspage.Statuspage.
func (t twilio) ScrapIncidents() ([]incidents.Incident, error) {
	var req atlassian.IncidentReq
	_, err := t.RestyClient.
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
func (t twilio) ScrapScheduleMaintainance() ([]incidents.Incident, error) {
	var req atlassian.IncidentReq
	_, err := t.RestyClient.
		R().
		SetResult(&req).
		Get(scheduleMaintenanceUrl)
	if err != nil {
		return nil, err
	}

	incidents := atlassian.FetchIncidentsHelper(req)

	return incidents, nil
}

func (t twilio) NewWithServiceID(id uint) statuspage.StatusPageProvider {
	t.ServiceID = id
	return t
}

func Register(client *resty.Client) {
	registry.Register(slug, twilio{
		RestyClient: client,
	})
}
