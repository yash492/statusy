package cmd

import (
	"backend/types"
	"testing"
)

func TestAtlassianScrapper(t *testing.T) {
	a := AtlassianIncidents{
		IncidentUrl: "https://status.circleci.com/api/v2/incidents.json",
		Incidents:   types.AtlassianStatusPageReq{},
	}

	a.ScrapIncidents()
}
