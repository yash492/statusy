package cmd

import (
	"backend/store"
	"backend/types"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gosimple/slug"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

type Services struct {
	Store store.ServicesStore
}

type Components struct {
	Store store.ComponentsStore
}

func AddServicesToDb(db *gorm.DB) error {
	bytes, err := os.ReadFile("./services.yaml")

	if err != nil {
		return fmt.Errorf(err.Error(), "could not read the services.yaml file")
	}

	var parseServices []types.Service
	err = yaml.Unmarshal(bytes, &parseServices)

	if err != nil {
		return fmt.Errorf(err.Error(), "could not unmarshal provider_details.yaml")
	}

	services := Services{
		Store: store.InitDbEnv(db),
	}

	_, err = services.Store.AddServices(parseServices)
	return err
}

func AddComponentsToDb(db *gorm.DB) error {

	services := Services{
		Store: store.InitDbEnv(db),
	}

	components := Components{
		Store: store.InitDbEnv(db),
	}

	fetchedServices, err := services.Store.GetAllServices()

	if err != nil {
		return err
	}

	totalComponents := make([]types.Component, 0)
	for _, service := range fetchedServices {
		if service.ProviderType == types.AtlassianProviderType {
			fetchedComponents, err := getAtlassianComponents(service.ComponentsUrl, service.ID)
			if err != nil {
				return err
			}

			totalComponents = append(totalComponents, fetchedComponents...)
		}
	}

	_, err = components.Store.AddComponents(totalComponents)
	return err
}

func getAtlassianComponents(url string, serviceId uint) ([]types.Component, error) {
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

	components := []types.Component{}

	for _, component := range atlassianComponentsReq.Components {
		components = append(components, types.Component{
			Name:      component.Name,
			ServiceId: serviceId,
			Slug:      slug.Make(component.Name),
		})
	}
	return components, nil
}
