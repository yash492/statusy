package scrapper

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"
	"unicode"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"github.com/gosimple/slug"
	"github.com/yash492/statusy/pkg/queue"
	"github.com/yash492/statusy/pkg/schema"
	"github.com/yash492/statusy/pkg/types"
)

type statusioProvider struct {
	incidentUrl string
	serviceID   uint
}

// Needs to re done
// Since this no longer is a valid approach
const incidentComponentsHeader = "Components"
const plannedMaintenance = "Planned Maintenance"
const breakDelimiter = "<br/>"

func (si statusioProvider) scrap(client *resty.Client, queue *queue.Queue) error {

	resp, err := client.R().Get(si.incidentUrl)
	if err != nil {
		return err
	}

	byteReader := bytes.NewReader(resp.Body())
	doc, err := goquery.NewDocumentFromReader(byteReader)
	if err != nil {
		return err
	}

	urlParsed, err := url.Parse(si.incidentUrl)
	if err != nil {
		return err
	}

	statusPageIncidents := make([]StatusPageIncident, 0)

	doc.Find(".timelineMinor").Each(func(i int, s *goquery.Selection) {
		// This could be incident status or scheduled maintenance
		eventStatus := s.Find(".pull-right.status_description").Text()
		if eventStatus == plannedMaintenance || eventStatus == "" {
			return
		}

		incidentTitleNode := s.Find(".panel-title a")
		incidentTitle := strings.TrimSpace(incidentTitleNode.Text())
		incidentLink, _ := incidentTitleNode.Attr("href")
		incidentComponentNames := make([]string, 0)
		incidentUpdates := make([]schema.IncidentUpdate, 0)

		if err != nil {
			return
		}

		// If incoming incident is after the latest stored incident for the given service then proceed else
		// If the last updated stored incident status time is greater than the incoming one then block the flow of the execution

		incidentLinkPath, _ := url.JoinPath(fmt.Sprintf("%v://%v", urlParsed.Scheme, urlParsed.Host), strings.TrimSpace(incidentLink))

		idxOfIncidentProviderIdFromPath := strings.LastIndex(incidentLink, "/")
		incidentProviderId := ""
		if idxOfIncidentProviderIdFromPath != -1 {
			incidentProviderId = incidentLink[idxOfIncidentProviderIdFromPath+1:]
		}

		if err != nil {
			log.Println(err)
		}

		s.Find(".panel > .panel-body").Children().Each(func(i int, s *goquery.Selection) {
			// This is incident status, components
			incidentEntity := strings.TrimSpace(s.Find(".event_inner_title").Text())

			// This is the value of the incident status, components
			incidentEntityValue := strings.TrimSpace(s.Find(".event_inner_text").Text())
			if incidentEntity == incidentComponentsHeader {
				incidentComponentNames = si.trimSpacesFromStringSlice(strings.Split(strings.TrimSpace(incidentEntityValue), ","))

			}

			// Incident Updates
			incidentTimeHtml, err := s.Find(".incident_time").Html()
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(incidentTimeHtml)

			incidentTimeStr := ""
			if strings.Contains(incidentTimeHtml, breakDelimiter) {
				incidentTimeStr = strings.SplitN(incidentTimeHtml, breakDelimiter, 2)[0]
			}

			var incidentStatusTime time.Time
			if incidentTimeStr != "" {
				incidentStatusTime, err = si.getTimeInUtc(incidentTimeStr)
				if err != nil {
					fmt.Println(err)
				}
			}

			incidentMsgSelection := s.Find(".incident_message_details")
			incidentUpdateMsgProviderID, _ := incidentMsgSelection.Attr("id")
			incidentUpdateMsg := incidentMsgSelection.Text()
			incidentUpdateProviderStatus := si.formatIncidentMsgStatus(incidentMsgSelection.Siblings().Text())

			if !incidentStatusTime.IsZero() {

				status := si.normaliseProviderState(incidentUpdateProviderStatus)
				// If it's the first incident status update we want the
				// normalised statusy status to triggered
				if i == 0 && incidentUpdateProviderStatus != "Resolved" {
					status = types.IncidentTriggered
				}

				incidentUpdates = append(incidentUpdates, schema.IncidentUpdate{
					Description:    incidentUpdateMsg,
					Status:         status,
					StatusTime:     incidentStatusTime,
					ProviderID:     incidentUpdateMsgProviderID,
					ProviderStatus: incidentUpdateProviderStatus,
				})
			}
		})

		if len(incidentUpdates) > 0 {
			incidentComponents := si.formatComponentNamesToStatusPageComponents(incidentComponentNames)

			impact := sql.NullString{
				String: eventStatus,
				Valid:  true,
			}

			incident := schema.Incident{
				Name:              incidentTitle,
				Link:              incidentLinkPath,
				ServiceID:         si.serviceID,
				ProviderID:        incidentProviderId,
				ProviderCreatedAt: incidentUpdates[0].StatusTime,
				ProviderImpact:    impact,
				Impact:            impact,
			}

			statusPageIncidents = append(statusPageIncidents, StatusPageIncident{
				Incident:           incident,
				IncidentUpdates:    incidentUpdates,
				IncidentComponents: incidentComponents,
			})
		}

	})

	return nil
}

func (statusioProvider) formatComponentNamesToStatusPageComponents(componentNames []string) []StatusPageIncidentComponent {
	components := make([]StatusPageIncidentComponent, 0)
	for _, componentName := range componentNames {
		components = append(components, StatusPageIncidentComponent{
			ComponentName: componentName,
			// Creating a slug for the component name
			// since status io doesn't provide any component id
			ProviderComponentID: slug.Make(componentName),
		})
	}
	return components
}

func (si statusioProvider) getTimeInUtc(incidentCreatedAtStr string) (time.Time, error) {
	daytime := "AM"
	if strings.Contains(incidentCreatedAtStr, "PM") {
		daytime = "PM"
	}

	index := strings.Index(incidentCreatedAtStr, daytime)
	incidentCreatedAtWithoutTz := strings.TrimSpace(incidentCreatedAtStr[:index+len(daytime)])
	incidentCreatedAtTz := strings.TrimSpace(incidentCreatedAtStr[index+len(daytime):])

	tz, err := si.getTimezoneForTzAbbr(incidentCreatedAtTz)
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

func (statusioProvider) normaliseProviderState(providerState string) string {
	stateMap := map[string]string{
		"Investigating": types.IncidentInProgress,
		"Identified":    types.IncidentInProgress,
		"Monitoring":    types.IncidentInProgress,
		"Resolved":      types.IncidentResolved,
	}
	return stateMap[providerState]
}

func (statusioProvider) trimSpacesFromStringSlice(slice []string) []string {
	trimmedSlice := make([]string, 0)
	for _, s := range slice {
		if strings.TrimSpace(s) == "" {
			continue
		}
		trimmedSlice = append(trimmedSlice, strings.TrimSpace(s))
	}
	return trimmedSlice
}

func (statusioProvider) formatIncidentMsgStatus(str string) string {
	trimmedStr := strings.TrimFunc(str, func(r rune) bool {
		return r == '[' || r == ']' || unicode.IsSpace(r)
	})

	return strings.ToLower(trimmedStr)

}

var statusioTimezoneAbbrMap = map[string]string{
	"SST":   "Pacific/Midway",
	"HST":   "Pacific/Honolulu",
	"AKST":  "America/Anchorage",
	"PST":   "America/Los_Angeles",
	"PDT":   "America/Los_Angeles",
	"MST":   "America/Denver",
	"CST":   "America/Chicago",
	"EST":   "America/New_York",
	"-04":   "America/Caracas",
	"-03":   "America/Santiago",
	"-02":   "America/Noronha",
	"-01":   "Atlantic/Cape_Verde",
	"UTC":   "UTC",
	"GMT":   "Europe/London",
	"CET":   "Europe/Brussels",
	"EET":   "Africa/Cairo",
	"SAST":  "Africa/Johannesburg",
	"MSK":   "Europe/Moscow",
	"+0330": "Asia/Tehran",
	"+04":   "Asia/Baku",
	"+0430": "Asia/Kabul",
	"+05":   "Asia/Tashkent",
	"IST":   "Asia/Kolkata",
	"+0530": "Asia/Colombo",
	"+0545": "Asia/Kathmandu",
	"+0630": "Indian/Cocos",
	"+07":   "Asia/Bangkok",
	"+08":   "Asia/Singapore",
	"KST":   "Asia/Seoul",
	"JST":   "Asia/Tokyo",
	"ACST":  "Australia/Darwin",
	"ChST":  "Pacific/Guam",
	"AEST":  "Australia/Brisbane",
	"AEDT":  "Australia/Melbourne",
	"+11":   "Australia/Lord_Howe",
	"+12":   "Pacific/Fiji",
	"NZDT":  "Pacific/Auckland",
	"+1345": "Pacific/Chatham",
	"+13":   "Pacific/Tongatapu",
}

func (statusioProvider) getTimezoneForTzAbbr(timezoneAbbr string) (string, error) {

	offset, ok := statusioTimezoneAbbrMap[timezoneAbbr]
	if !ok {
		return "", fmt.Errorf("can find timezone name for the given abbr %v", timezoneAbbr)
	}
	return offset, nil
}
