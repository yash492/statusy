package cmd

import (
	"backend/config"
	"backend/store"
	"backend/types"
	"testing"
)

func TestAtlassianScrapper(t *testing.T) {
	config.Load()
	store.InitDb()
	initRepos()
	a := AtlassianIncidents{
		IncidentUrl:  "https://status.circleci.com/api/v2/incidents.json",
		Incidents:    types.AtlassianStatusPageReq{},
		ProviderSlug: "circleci",
	}

	a.ScrapIncidents()
}
