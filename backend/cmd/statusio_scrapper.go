package cmd

import (
	"backend/external"
	"backend/models"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
	_ "time/tzdata"
	"unicode"

	"github.com/PuerkitoBio/goquery"
)

type StatusioIncidents struct {
	IncidentUrl  string
	Incidents    []models.Incident
	ProviderSlug string
}

const incidentStatusHeader = "Incident Status"
const incidentComponentsHeader = "Components"
const plannedMaintenance = "Planned Maintenance"
const breakDelimiter = "<br/>"

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
	incident := &models.Incident{}

	urlParsed, err := url.Parse(s.IncidentUrl)
	if err != nil {
		return err
	}

	doc.Find(".timelineMinor").Each(func(i int, s *goquery.Selection) {
		// This could be incident status or scheduled maintenance
		eventStatus := s.Find(".pull-right.status_description").Text()
		if eventStatus == plannedMaintenance || eventStatus == "" {
			return
		}

		incidentTitleNode := s.Find(".panel-title a")
		incidentTitle := trim(incidentTitleNode.Text())
		incidentLink, _ := incidentTitleNode.Attr("href")

		incidentComponents := make([]string, 0)
		incidentStatus := ""

		s.Find(".panel > .panel-body").Children().Each(func(i int, s *goquery.Selection) {
			// This is incident status, components
			incidentEntity := trim(s.Find(".event_inner_title").Text())

			// This is value of the incident status, components
			incidentEntityValue := trim(s.Find(".event_inner_text").Text())
			if incidentEntity == incidentStatusHeader {
				incidentStatus = incidentEntityValue
			}

			if incidentEntity == incidentComponentsHeader {
				incidentComponents = trimSpacesFromStringSlice(strings.Split(trim(incidentEntityValue), ","))
			}

			// Incident Updates
			incidentTimeHtml, err := s.Find(".incident_time").Html()
			if err != nil {
				fmt.Println(err)
			}
			incidentTimeStr := ""
			if strings.Contains(incidentTimeHtml, breakDelimiter) {
				incidentTimeStr = strings.SplitN(incidentTimeHtml, breakDelimiter, 2)[0]
			}

			var incidentTime time.Time
			if incidentTimeStr != "" {
				incidentTime, _ = getTimeInUtc(incidentTimeStr)
			}

			incidentMsgSelection := s.Find(".incident_message_details")
			incidentMsg := incidentMsgSelection.Text()
			incidentMsgStatus := formatIncidentMsgStatus(incidentMsgSelection.Siblings().Text())

			if !incidentTime.IsZero() {
				fmt.Println(incidentMsg, incidentMsgStatus, incidentTime)
			}

		})

		incident.Description = incidentTitle
		joinedUrlPath, _ := url.JoinPath(fmt.Sprintf("%v://%v", urlParsed.Scheme, urlParsed.Host), trim(incidentLink))
		incident.Url = joinedUrlPath

		idxOfIncidentProviderIncIdFromPath := strings.LastIndex(incidentLink, "/")
		if idxOfIncidentProviderIncIdFromPath != -1 {
			incident.ProviderIncidentId = incidentLink[idxOfIncidentProviderIncIdFromPath+1:]
		}

		fmt.Println(incidentComponents, incidentStatus)

	})

	return nil
}

func getTimeInUtc(incidentCreatedAtStr string) (time.Time, error) {
	daytime := "AM"
	if strings.Contains(incidentCreatedAtStr, "PM") {
		daytime = "PM"
	}

	index := strings.Index(incidentCreatedAtStr, daytime)
	incidentCreatedAtWithoutTz := trim(incidentCreatedAtStr[:index+len(daytime)])
	incidentCreatedAtTz := trim(incidentCreatedAtStr[index+len(daytime):])

	tz, err := external.GetTimezoneForTzAbbr(incidentCreatedAtTz)
	if err != nil {
		return time.Time{}, err
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		return time.Time{}, err
	}

	parsedTime, err := time.ParseInLocation("January 2, 2006 3:04PM", incidentCreatedAtWithoutTz, loc)
	return parsedTime.UTC(), err

}

func trimSpacesFromStringSlice(slice []string) []string {
	trimmedSlice := make([]string, 0)
	for _, s := range slice {
		if trim(s) == "" {
			continue
		}
		trimmedSlice = append(trimmedSlice, trim(s))
	}

	return trimmedSlice
}

func formatIncidentMsgStatus(str string) string {
	trimmedStr := strings.TrimFunc(str, func(r rune) bool {
		return r == '[' || r == ']' || unicode.IsSpace(r)
	})

	return strings.ToLower(trimmedStr)

}
