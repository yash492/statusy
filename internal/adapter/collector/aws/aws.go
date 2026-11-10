package aws

import (
	"github.com/yash492/statusy/internal/adapter/collector/registry"
	"github.com/yash492/statusy/internal/domain/components"
	"github.com/yash492/statusy/internal/domain/incidents"
	"github.com/yash492/statusy/internal/domain/scheduledmaintenance"
	"github.com/yash492/statusy/internal/domain/services"
	"github.com/yash492/statusy/internal/domain/statuspage"
	"resty.dev/v3"
)

const slug = "aws"
const name = "Amazon Web Services (AWS)"

const (
	historyEventsUrl = "https://history-events-ap-northeast-1-prod.s3.amazonaws.com/historyevents.json"
	servicesUrl      = "https://servicedata-ap-northeast-1-prod.s3.amazonaws.com/services.json"
	statusPageUrl    = "https://health.aws.amazon.com/health/status"
)

type aws struct {
	RestyClient *resty.Client
	ServiceID   uint
}

func (a aws) ID() uint {
	return a.ServiceID
}

func (a aws) Name() string {
	return name
}

var _ statuspage.StatusPageProvider = aws{}

// ScrapComponents implements statuspage.StatusPageProvider.
func (a aws) ScrapComponents() (components.AggregateComponents, error) {
	var req ServicesReq
	_, err := a.RestyClient.
		R().
		SetResult(&req).
		EnableRetryDefaultConditions().
		Get(servicesUrl)
	if err != nil {
		return components.AggregateComponents{}, err
	}

	return fetchComponentsHelper(req), nil
}

// Slug implements statuspage.StatusPageProvider.
func (a aws) Slug() services.ServiceSlug {
	return slug
}

// ScrapIncidents implements statuspage.StatusPageProvider.
func (a aws) ScrapIncidents() ([]incidents.Incident, error) {
	var req HistoryEventsReq
	_, err := a.RestyClient.
		R().
		SetResult(&req).
		Get(historyEventsUrl)
	if err != nil {
		return nil, err
	}

	return fetchIncidentsHelper(req), nil
}

// ScrapScheduledMaintenance implements statuspage.StatusPageProvider.
// AWS does not publish scheduled maintenance on its public dashboard.
func (a aws) ScrapScheduledMaintenance() ([]scheduledmaintenance.ScheduledMaintenance, error) {
	return []scheduledmaintenance.ScheduledMaintenance{}, nil
}

func (a aws) GetStatuspageUrl() string {
	return statusPageUrl
}

func (a aws) NewWithServiceID(id uint) statuspage.StatusPageProvider {
	a.ServiceID = id
	return a
}

func Register(client *resty.Client) {
	registry.Register(slug, aws{
		RestyClient: client,
	})
}
