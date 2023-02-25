package cmd

import (
	"backend/models"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type StatusioIncidents struct {
	IncidentUrl  string
	Incidents    []models.Incident
	ProviderSlug string
}

var trim = strings.TrimSpace

func (s *StatusioIncidents) ScrapIncidents() error {

	resp, err := http.Get(s.IncidentUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	doc.Find("#statusio_history_timeline").Each(func(i int, s *goquery.Selection) {
		s.Find(".history_timeline_event_title").Each(func(i int, s *goquery.Selection) {
			// incidentTitleParsed := s.Find("a").Text()
			// incidentCreatedAtStrParsed := s.Find(".history_timeline_event_date").Text()

			// fmt.Println(incidentCreatedAtStrParsed, incidentTitleParsed)
		})

		s.Find(".panel > .panel-body").Each(func(i int, s *goquery.Selection) {
			
			// Using hyphen as a delimiter to fetch components, incident status
			incidentMetadataHeader := s.Find(".event_inner_title").AppendHtml("<p>-</p>").Text()
			incidentMetadataValue := s.Find(".event_inner_text").AppendHtml("<p>-</p>").Text()

			mappedIncidentMetadata := mapIncidentMetadata(strings.Split(incidentMetadataHeader, "-"), strings.Split(incidentMetadataValue, "-"))
			fmt.Println(mappedIncidentMetadata)
		})

	})

	return nil
}

func mapIncidentMetadata(headers []string, values []string) map[string]string {

	metaMap := make(map[string]string, len(headers))
	for i, header := range headers {
		if header == "" || values[i] == "" {
			continue
		}
		metaMap[trim(header)] = trim(values[i])
	}

	return metaMap
}
