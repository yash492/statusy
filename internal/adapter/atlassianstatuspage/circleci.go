package atlassianstatuspage

import (
	"encoding/json"

	"github.com/yash492/statusy/internal/domain/statuspage"
	"resty.dev/v3"
)

const circleciSlug = "circleSlug"

type CircleCi struct {
	ScheduleMaintainanceUrl string
	IncidentsUrl            string
	ComponentsUrl           string
	RestyClient             *resty.Client
}

func New(
	componentsUrl string,
	incidentsUrl string,
	scheduleMaintainanceUrl string,
	restyClient *resty.Client,
) CircleCi {
	return CircleCi{
		IncidentsUrl:            incidentsUrl,
		ComponentsUrl:           componentsUrl,
		ScheduleMaintainanceUrl: scheduleMaintainanceUrl,
		RestyClient:             restyClient,
	}
}

var _ statuspage.StatusPage = CircleCi{}

// FetchComponents implements statuspage.Statuspage.
func (c CircleCi) FetchComponents() (statuspage.ComponentGroups, error) {
	resp, err := c.RestyClient.R().Get(c.ComponentsUrl)
	if err != nil {
		return nil, err
	}

	circleciComponents := atlassianComponentsReq{}
	err = json.Unmarshal(resp.Bytes(), &circleciComponents)
	if err != nil {
		return nil, err
	}

	

}

// Name implements statuspage.Statuspage.
func (c CircleCi) Slug() string {
	return circleciSlug
}

// ScrapIncidents implements statuspage.Statuspage.
func (c CircleCi) ScrapIncidents() {
	panic("unimplemented")
}

// ScrapScheduleMaintainance implements statuspage.Statuspage.
func (c CircleCi) ScrapScheduleMaintainance() {
	panic("unimplemented")
}
