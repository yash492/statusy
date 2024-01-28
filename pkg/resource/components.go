package resource

import (
	"bytes"
	"log/slog"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"github.com/gosimple/slug"
	"github.com/samber/lo"
	"github.com/yash492/statusy/pkg/domain"
	"github.com/yash492/statusy/pkg/schema"
	"github.com/yash492/statusy/pkg/types"
)

type atlassianComponentsReq struct {
	Components []atlassianComponent `json:"components"`
}

type atlassianComponent struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	Status             string    `json:"status"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	Position           int       `json:"position"`
	Description        *string   `json:"description"`
	Showcase           bool      `json:"showcase"`
	GroupID            string    `json:"group_id"`
	PageID             string    `json:"page_id"`
	Group              bool      `json:"group"`
	OnlyShowIfDegraded bool      `json:"only_show_if_degraded"`
}

type parseProviderComponentsFunc func(client *resty.Client, serviceId uint, componentsUrl string) ([]schema.Component, error)

var componentsParseMap = map[string]parseProviderComponentsFunc{
	types.AtlassianProviderType: parseAtlassianComponents,
	types.StatusioProviderType:  parseStatusioComponents,
}

func initComponents() error {
	services, err := domain.Service.GetAll()
	if err != nil {
		return err
	}

	restyClient := resty.New()

	var components []schema.Component

	for _, service := range services {
		if parseComponentFunc, ok := componentsParseMap[service.ProviderType]; ok {
			parsedComponents, err := parseComponentFunc(restyClient, service.ID, service.ComponentsURL)
			if err != nil {
				return err
			}
			components = append(components, parsedComponents...)
		}
	}

	_, err = domain.Component.Create(components)
	if err != nil {
		return err
	}
	return nil
}

func parseAtlassianComponents(
	client *resty.Client,
	serviceId uint,
	componentsUrl string) ([]schema.Component, error) {

	var atlassianComponentReq atlassianComponentsReq
	_, err := client.R().SetResult(&atlassianComponentReq).Get(componentsUrl)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	var components []schema.Component
	for _, componentReq := range atlassianComponentReq.Components {

		if componentReq.Group {
			continue
		}

		components = append(components, schema.Component{
			Name:       componentReq.Name,
			ServiceID:  serviceId,
			ProviderID: componentReq.ID,
		})
	}

	return components, nil

}

func parseStatusioComponents(
	client *resty.Client,
	serviceId uint,
	componentsUrl string) ([]schema.Component, error) {

	resp, err := client.R().Get(componentsUrl)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body()))
	if err != nil {
		return nil, err
	}

	componentNames := make([]string, 0)
	doc.Find("#statusio_components .component_name").Each(func(i int, s *goquery.Selection) {
		componentNames = append(componentNames, strings.TrimSpace(s.Text()))
	})

	components := lo.Map(componentNames, func(componentName string, _ int) schema.Component {
		return schema.Component{
			Name:       componentName,
			ProviderID: slug.Make(componentName),
			ServiceID:  serviceId,
		}
	})

	return components, nil

}
