package resource

import (
	"log/slog"

	"github.com/go-resty/resty/v2"
	"github.com/yash492/statusy/pkg/domain"
	"github.com/yash492/statusy/pkg/schema"
	"github.com/yash492/statusy/pkg/types"
)

type parseProviderComponentsFunc func(client *resty.Client, serviceId uint, componentsUrl string) ([]schema.Component, error)

var componentsParseMap = map[string]parseProviderComponentsFunc{
	types.AtlassianProviderType: parseAtlassianComponents,
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

	err = domain.Component.Create(components)
	if err != nil {
		return err
	}
	return nil
}

func parseAtlassianComponents(
	client *resty.Client,
	serviceId uint,
	componentsUrl string) ([]schema.Component, error) {

	var atlassianComponentReq types.AtlassianComponentsReq
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
			Name:                componentReq.Name,
			ServiceId:           serviceId,
			ProviderComponentId: componentReq.ID,
		})
	}

	return components, nil

}
