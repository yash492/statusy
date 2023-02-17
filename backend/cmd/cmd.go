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

type ServicesEnv struct {
	Store store.ServicesStore
}

type ComponentsEnv struct {
	Store store.ComponentsStore
}

type IncidentsEnv struct {
	Store store.IncidentStore
}

type IncidentUpdatesEnv struct {
	Store store.IncidentUpdateStore
}

type IncidentComponentsEnv struct {
	Store store.IncidentComponentsStore
}

var Components = ComponentsEnv{
	Store: store.InitDbEnv(),
}

var Services = ServicesEnv{
	Store: store.InitDbEnv(),
}

var Incidents = IncidentsEnv{
	Store: store.InitDbEnv(),
}

var IncidentUpdate = IncidentUpdatesEnv{
	Store: store.InitDbEnv(),
}

var IncidentUpdatesComponent = IncidentComponentsEnv{
	Store: store.InitDbEnv(),
}

func AddServicesAndComponentsToDb(db *gorm.DB) error {
	bytes, err := os.ReadFile("./services.yaml")
	if err != nil {
		return fmt.Errorf(err.Error(), "could not read the services.yaml file")
	}

	var parseServices []types.Service

	err = yaml.Unmarshal(bytes, &parseServices)
	if err != nil {
		return fmt.Errorf(err.Error(), "could not unmarshal provider_details.yaml")
	}

	insertedServices, err := Services.Store.AddServices(parseServices)
	if err != nil {
		return err
	}

	totalComponents := make([]types.Component, 0)
	for _, service := range insertedServices {
		if service.ProviderType == types.AtlassianProviderType {
			fetchedComponents, err := getAtlassianComponents(service.ComponentsUrl, service.ID)
			if err != nil {
				return err
			}
			totalComponents = append(totalComponents, fetchedComponents...)
		}
	}

	_, err = Components.Store.AddComponents(totalComponents)
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
