package cmd

import (
	"backend/config"
	"backend/models"
	"backend/store"
	"backend/types"
	"testing"
)

func TestAtlassianScrapper(t *testing.T) {
	config.Load("../config.yaml")
	store.InitDb()
	initRepos()
	a := AtlassianIncidents{
		IncidentUrl:  "https://status.circleci.com/api/v2/incidents.json",
		Incidents:    types.AtlassianStatusPageReq{},
		ProviderSlug: "circleci",
	}

	err := a.ScrapIncidents()
	if err != nil {
		panic(err)
	}
}

func TestStatusioScrapper(t *testing.T) {
	config.Load("../config.yaml")
	store.InitDb()
	initRepos()
	a := StatusioIncidents{
		IncidentUrl:  "https://status.roblox.com/pages/history/59db90dbcdeb2f04dadcf16d",
		Incidents:    []models.Incident{},
		ProviderSlug: "docker",
	}

	err := a.ScrapIncidents()
	if err != nil {
		panic(err)
	}
}
