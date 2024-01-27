package scrapper

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"github.com/yash492/statusy/pkg/queue"
	"github.com/yash492/statusy/pkg/schema"
	"github.com/yash492/statusy/pkg/types"
)

type squadcastProvider struct {
	incidentUrl string
	serviceID   uint
}

// Scraping from Squadcast web page itselt
// as there isn't any proper api for this yet
type squadcastReq struct {
	Props propsReq `json:"props"`
}

type propsReq struct {
	PageProps pagePropsReq `json:"pageProps"`
}

type pagePropsReq struct {
	History []historyReq `json:"history"`
}

type historyReq struct {
	Date   time.Time   `json:"date"`
	Issues []issuesReq `json:"issues"`
}

// this corresponds to incident
type issuesReq struct {
	ID           json.Number `json:"id"` // incident provider_id
	Title        string      `json:"title"`
	CreatedAt    time.Time   `json:"createdAt"`
	ResolvedAt   time.Time   `json:"resolvedAt"`
	CurrentState string      `json:"currentState"`
	States       []states    `json:"states"`
}

// this corresponds to incident updates
type states struct {
	// Incident Update Status
	Name     string `json:"name"`
	Messages []struct {
		ID        json.Number `json:"id"` //incident_update provider id
		Text      string      `json:"text"`
		Timestamp time.Time   `json:"timestamp"`
	} `json:"messages"`
}

func (s squadcastProvider) scrap(client *resty.Client, queue *queue.Queue) error {
	resp, err := client.R().Execute(resty.MethodGet, s.incidentUrl)
	if err != nil {
		return err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body()))
	if err != nil {
		return err
	}
	selection := doc.Find("script#__NEXT_DATA__")

	var squadcastStatusPage squadcastReq
	err = json.Unmarshal([]byte(selection.Text()), &squadcastStatusPage)
	if err != nil {
		return err
	}

	incidents := s.normaliseIncidents(squadcastStatusPage.Props.PageProps.History)

	return storeAndDispatchIncidents(queue, incidents)
}

// normaliseIncidents normalises custom statuspage response to
// the statusy's schema
// for example: squadcastProvider -> StatusPageIncident
func (s squadcastProvider) normaliseIncidents(histories []historyReq) []StatusPageIncident {

	var statusPageIncidents []StatusPageIncident
	squadcastIncidentMaps := make(map[string]StatusPageIncident, 0)

	for _, history := range histories {
		if len(history.Issues) < 1 {
			continue
		}

		for _, issue := range history.Issues {
			incident := schema.Incident{
				Name:              issue.Title,
				ProviderID:        issue.ID.String(),
				ProviderCreatedAt: issue.CreatedAt,
				ServiceID:         s.serviceID,
				Link:              s.incidentUrl,
			}

			var incidentUpdate []schema.IncidentUpdate
			for _, state := range issue.States {
				incidentProviderStatus := state.Name
				for i, message := range state.Messages {

					status := normaliseProviderState(incidentProviderStatus)
					// If it's the first incident status update we want the
					// normalised statusy status to triggered
					if i == 0 && incidentProviderStatus == "Investigating" {
						status = types.IncidentTriggered
					}

					incidentUpdate = append(incidentUpdate, schema.IncidentUpdate{
						Description:    message.Text,
						ProviderID:     message.ID.String(),
						StatusTime:     message.Timestamp,
						Status:         status,
						ProviderStatus: incidentProviderStatus,
					})
				}
			}

			// If the incident is already present in the map
			// we don't want to add it again
			if _, ok := squadcastIncidentMaps[incident.ProviderID]; !ok {
				statusPageIncident := StatusPageIncident{
					Incident:        incident,
					IncidentUpdates: incidentUpdate,
				}

				squadcastIncidentMaps[incident.ProviderID] = statusPageIncident
				statusPageIncidents = append(statusPageIncidents, statusPageIncident)
			}
		}
	}

	return statusPageIncidents
}

func normaliseProviderState(providerState string) string {
	stateMap := map[string]string{
		"Investigating": types.IncidentInProgress,
		"Identified":    types.IncidentInProgress,
		"Monitoring":    types.IncidentInProgress,
		"Resolved":      types.IncidentResolved,
	}
	return stateMap[providerState]
}
