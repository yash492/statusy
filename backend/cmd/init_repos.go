package cmd

import "backend/store"

type servicesConfig struct {
	Store store.ServicesStore
}

type componentsConfig struct {
	Store store.ComponentsStore
}

type incidentsConfig struct {
	Store store.IncidentStore
}

type incidentUpdatesConfig struct {
	Store store.IncidentUpdateStore
}

type incidentComponentsConfig struct {
	Store store.IncidentComponentsStore
}

var componentsEnv componentsConfig
var incidentsEnv incidentsConfig
var servicesEnv servicesConfig
var incidentUpdatesEnv incidentUpdatesConfig
var incidentComponentsEnv incidentComponentsConfig

func initRepos() {
	componentsEnv = componentsConfig{Store: store.InitDbEnv()}
	servicesEnv = servicesConfig{Store: store.InitDbEnv()}
	incidentsEnv = incidentsConfig{Store: store.InitDbEnv()}
	incidentComponentsEnv = incidentComponentsConfig{Store: store.InitDbEnv()}
	incidentUpdatesEnv = incidentUpdatesConfig{Store: store.InitDbEnv()}
}
