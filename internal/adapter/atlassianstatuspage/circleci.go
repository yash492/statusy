package atlassianstatuspage

import (
	"encoding/json"

	"github.com/samber/lo"
	"github.com/yash492/statusy/internal/domain/statuspage"
	"resty.dev/v3"
)

const circleciSlug = "circleci"

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
func (c CircleCi) FetchComponents() (statuspage.AggregateComponents, error) {
	resp, err := c.RestyClient.R().Get(c.ComponentsUrl)
	if err != nil {
		return statuspage.AggregateComponents{}, err
	}

	circleciComponents := atlassianComponentsReq{}
	err = json.Unmarshal(resp.Bytes(), &circleciComponents)
	if err != nil {
		return statuspage.AggregateComponents{}, err
	}

	componentGroups := c.fetchComponentsHelper(circleciComponents)
	return componentGroups, nil
}

func (c CircleCi) fetchComponentsHelper(atlassianComponents atlassianComponentsReq) statuspage.AggregateComponents {

	ungroupedComponents := []statuspage.Component{}
	groupedComponentsIDNameMap := map[string]string{}
	componentsToBeGrouped := []atlassianComponent{}
	groupedComponents := map[string]statuspage.ComponentGroup{}

	for _, atlassianComponent := range atlassianComponents.Components {
		if atlassianComponent.GroupID == nil {
			if !atlassianComponent.Group {
				ungroupedComponents = append(ungroupedComponents, statuspage.Component{
					Name:        atlassianComponent.Name,
					ServiceSlug: c.Slug(),
					ProviderID:  atlassianComponent.ID,
				})

			} else if atlassianComponent.Group {
				groupedComponentsIDNameMap[atlassianComponent.ID] = atlassianComponent.Name
			}
		} else {
			componentsToBeGrouped = append(componentsToBeGrouped, atlassianComponent)
		}
	}

	for _, altassianComponent := range componentsToBeGrouped {
		if altassianComponent.GroupID == nil {
			continue
		}
		componentGroupName := groupedComponentsIDNameMap[*altassianComponent.GroupID]
		componentGroupID := *altassianComponent.GroupID
		groupedComponent, ok := groupedComponents[*altassianComponent.GroupID]
		if !ok {
			groupedComponent = statuspage.ComponentGroup{
				Name:       componentGroupName,
				ProviderID: componentGroupID,
				Components: []statuspage.Component{},
			}
		}
		component := statuspage.Component{
			Name:        altassianComponent.Name,
			ServiceSlug: c.Slug(),
			ProviderID:  altassianComponent.ID,
		}

		groupedComponent.Components = append(groupedComponent.Components, component)
		groupedComponents[*altassianComponent.GroupID] = groupedComponent

	}

	return statuspage.AggregateComponents{
		UngroupedComponents: ungroupedComponents,
		GroupedComponents:   lo.Values(groupedComponents),
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
