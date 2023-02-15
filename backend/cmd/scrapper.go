package cmd

import (
	"backend/types"
	"encoding/json"
	"io"
	"net/http"
)

type Scrapper interface {
	ScrapIncidents() error
	ScrapMaintenance() error
}

type AtlassianIncidents struct {
	IncidentUrl string
	Incidents   types.AtlassianStatusPageReq
}

func (a *AtlassianIncidents) ScrapIncidents() error {
	resp, err := http.Get(a.IncidentUrl)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &a.Incidents)
	if err != nil {
		return err
	}

	return nil
}
