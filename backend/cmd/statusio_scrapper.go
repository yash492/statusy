package cmd

import (
	"backend/models"
	"fmt"
	"net/http"
	"strings"
	"time"
	_ "time/tzdata"

	"github.com/PuerkitoBio/goquery"
)

type StatusioIncidents struct {
	IncidentUrl  string
	Incidents    []models.Incident
	ProviderSlug string
}

const incidentStatusHeader = "Incident Status"
const componentsHeader = "Conponents"

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
			incidentTitleParsed := s.Find("a").Text()
			incidentCreatedAtStrParsed := s.Find(".history_timeline_event_date").Text()

			fmt.Println(incidentCreatedAtStrParsed, incidentTitleParsed)
		})

		s.Find(".panel > .panel-body").Each(func(i int, s *goquery.Selection) {

			incidentMetadataHeaders := make([]string, 0)
			s.Find(".event_inner_title").Each(func(i int, s *goquery.Selection) {
				text := s.Text()
				if text == "" {
					return
				}

				incidentMetadataHeaders = append(incidentMetadataHeaders, trim(text))
			})

			incidentMetadataValues := make([]string, 0)
			s.Find(".event_inner_text").Each(func(i int, s *goquery.Selection) {
				text := s.Text()
				if text == "" {
					return
				}
				incidentMetadataValues = append(incidentMetadataValues, trim(text))
			})

			mappedIncidentMetadata := mapIncidentMetadata(incidentMetadataHeaders, incidentMetadataValues)
			fmt.Println(mappedIncidentMetadata)

			s.Find(".incident_time:first-child").Each(func(i int, s *goquery.Selection) {
				html, _ := s.Html()
				delimiter := "<br/>"
				incidentCreationTimeStr := ""
				if strings.Contains(html, delimiter) {
					incidentCreationTimeStr = strings.SplitN(html, delimiter, 2)[0]
				} else {
					incidentCreationTimeStr = html
				}
				loc, err := time.LoadLocation("America/Chicago")


				t, err := time.ParseInLocation("January 2, 2006 3:04PM MST", incidentCreationTimeStr, loc)
				fmt.Println(t.String(), incidentCreationTimeStr, err)

			})
		})

	})

	return nil
}

func mapIncidentMetadata(headers []string, values []string) map[string]string {

	incidentMetadataMap := make(map[string]string, len(headers))
	for i, header := range headers {
		incidentMetadataMap[header] = values[i]
	}

	return incidentMetadataMap
}
