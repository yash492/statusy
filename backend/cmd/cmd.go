package cmd

import (
	"backend/models"
	"backend/types"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gosimple/slug"
	"gopkg.in/yaml.v3"
)

func Init() error {
	initRepos()

	err := AddServicesToDb()
	if err != nil {
		return err
	}

	err = AddComponentsToDb()
	if err != nil {
		return err
	}

	return nil
}

func AddServicesToDb() error {
	bytes, err := os.ReadFile("./services.yaml")
	if err != nil {
		return fmt.Errorf(err.Error(), "could not read the services.yaml file")
	}

	var parseServices []models.Service

	err = yaml.Unmarshal(bytes, &parseServices)
	if err != nil {
		return fmt.Errorf(err.Error(), "could not unmarshal provider_details.yaml")
	}

	_, err = servicesEnv.Store.AddServices(parseServices)
	return err
}

func AddComponentsToDb() error {

	services, err := servicesEnv.Store.GetAllServices()
	if err != nil {
		return err
	}

	totalComponents := make([]models.Component, 0)
	for _, service := range services {
		if service.ProviderType == types.AtlassianProviderType {
			fetchedComponents, err := getAtlassianComponents(service.ComponentsUrl, service.ID)
			if err != nil {
				return err
			}
			totalComponents = append(totalComponents, fetchedComponents...)
		}
	}

	_, err = componentsEnv.Store.AddComponents(totalComponents)
	return err
}

func getAtlassianComponents(url string, serviceId uint) ([]models.Component, error) {

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	atlassianComponentsReq := types.AtlassianComponentsReq{}
	err = json.Unmarshal(bytes, &atlassianComponentsReq)
	if err != nil {
		return nil, err
	}

	components := []models.Component{}
	for _, component := range atlassianComponentsReq.Components {
		components = append(components, models.Component{
			Name:      component.Name,
			ServiceId: serviceId,
			Slug:      slug.Make(component.Name),
		})
	}
	return components, nil
}
