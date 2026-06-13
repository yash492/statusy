package incidentio

import (
	"time"

	"github.com/yash492/statusy/internal/common/nullable"
	"github.com/yash492/statusy/internal/common/statusnormalizer"
	"github.com/yash492/statusy/internal/domain/components"
	"github.com/yash492/statusy/internal/domain/incidents"
	"github.com/yash492/statusy/internal/domain/scheduledmaintenance"
)

// StatusReq represents the statuspage summary JSON response.
// The Summary API is accessed via a GET request using the Resty client.
// The response JSON contains components and component structure items, which are parsed by the helper.
type StatusReq struct {
	Summary Summary `json:"summary"`
}

type Summary struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	PublicURL  string      `json:"public_url"`
	Components []Component `json:"components"`
	Structure  Structure   `json:"structure"`
}

type Component struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	StatusPageID string `json:"status_page_id"`
}

type Structure struct {
	ID           string          `json:"id"`
	StatusPageID string          `json:"status_page_id"`
	Items        []StructureItem `json:"items"`
}

type StructureItem struct {
	Group *Group `json:"group,omitempty"`
}

type Group struct {
	ID         string           `json:"id"`
	Name       string           `json:"name"`
	Components []GroupComponent `json:"components"`
}

type GroupComponent struct {
	ComponentID string `json:"component_id"`
	Name        string `json:"name"`
}

// IncidentsReq represents the incidents history JSON response.
type IncidentsReq struct {
	Incidents []Incident `json:"incidents"`
}

type Incident struct {
	ID                 string              `json:"id"`
	Name               string              `json:"name"`
	Status             string              `json:"status"`
	PublishedAt        time.Time           `json:"published_at"`
	Updates            []Update            `json:"updates"`
	AffectedComponents []AffectedComponent `json:"affected_components"`
	StatusSummaries    []StatusSummary     `json:"status_summaries"`
	Type               string              `json:"type"`
}

type Update struct {
	ID                string            `json:"id"`
	PublishedAt       time.Time         `json:"published_at"`
	MessageString     string            `json:"message_string"`
	ToStatus          string            `json:"to_status"`
	ComponentStatuses []ComponentStatus `json:"component_statuses"`
}

type ComponentStatus struct {
	ComponentID string `json:"component_id"`
	Status      string `json:"status"`
}

type AffectedComponent struct {
	ComponentID   string `json:"component_id"`
	Status        string `json:"status"`
	CurrentStatus string `json:"current_status"`
}

type StatusSummary struct {
	StartAt              time.Time  `json:"start_at"`
	EndAt                *time.Time `json:"end_at,omitempty"`
	WorstComponentStatus string     `json:"worst_component_status"`
}

// FetchComponentsHelper parses StatusReq, mapping groups and components to the domain models.
func FetchComponentsHelper(req StatusReq) components.AggregateComponents {
	groupedComponents := []components.ComponentGroup{}
	groupedComponentIDs := make(map[string]bool)

	for _, item := range req.Summary.Structure.Items {
		if item.Group != nil {
			g := components.ComponentGroup{
				Name:       item.Group.Name,
				ProviderID: item.Group.ID,
				Components: []components.Component{},
			}
			for _, c := range item.Group.Components {
				g.Components = append(g.Components, components.Component{
					Name:       c.Name,
					ProviderID: c.ComponentID,
				})
				groupedComponentIDs[c.ComponentID] = true
			}
			groupedComponents = append(groupedComponents, g)
		}
	}

	ungroupedComponents := []components.Component{}
	for _, c := range req.Summary.Components {
		if !groupedComponentIDs[c.ID] {
			ungroupedComponents = append(ungroupedComponents, components.Component{
				Name:       c.Name,
				ProviderID: c.ID,
			})
		}
	}

	return components.AggregateComponents{
		GroupedComponents:   groupedComponents,
		UngroupedComponents: ungroupedComponents,
	}
}

// FetchIncidentsHelper parses IncidentsReq, mapping historical incidents and updates to the domain models.
func FetchIncidentsHelper(req IncidentsReq, statusPageUrl string) []incidents.Incident {
	incidentList := []incidents.Incident{}
	for _, incidentReq := range req.Incidents {
		if incidentReq.Type != "" && incidentReq.Type != "incident" {
			continue
		}

		worstImpact := ""
		for _, summary := range incidentReq.StatusSummaries {
			worstImpact = getWorseStatus(worstImpact, summary.WorstComponentStatus)
		}

		incident := incidents.Incident{
			Name:              incidentReq.Name,
			Link:              statusPageUrl + "/incidents/" + incidentReq.ID,
			ProviderImpact:    nullable.SetValue(worstImpact, worstImpact != ""),
			Impact:            nullable.SetValue(worstImpact, worstImpact != ""),
			ProviderID:        incidentReq.ID,
			ProviderCreatedAt: incidentReq.PublishedAt,
			Components:        []components.Component{},
		}

		for _, comp := range incidentReq.AffectedComponents {
			incident.Components = append(incident.Components, components.Component{
				ProviderID: comp.ComponentID,
			})
		}

		incident.Updates = make([]incidents.IncidentUpdate, len(incidentReq.Updates))
		for idx, update := range incidentReq.Updates {
			incident.Updates[idx] = incidents.IncidentUpdate{
				Description:        update.MessageString,
				IncidentProviderID: incident.ProviderID,
				ProviderID:         update.ID,
				ProviderStatus:     update.ToStatus,
				Status:             statusnormalizer.NormalizeIncidentStatus(update.ToStatus),
				StatusTime:         update.PublishedAt,
			}
		}

		incidentList = append(incidentList, incident)
	}

	return incidentList
}

// FetchScheduledMaintenancesHelper parses IncidentsReq, filtering and mapping historical scheduled maintenance windows to domain models.
func FetchScheduledMaintenancesHelper(req IncidentsReq, statusPageUrl string) []scheduledmaintenance.ScheduledMaintenance {
	maintenanceList := []scheduledmaintenance.ScheduledMaintenance{}
	for _, incidentReq := range req.Incidents {
		if incidentReq.Type != "maintenance" {
			continue
		}

		var startsAt time.Time
		var endsAt time.Time

		if len(incidentReq.StatusSummaries) > 0 {
			summary := incidentReq.StatusSummaries[0]
			startsAt = summary.StartAt
			if summary.EndAt != nil {
				endsAt = *summary.EndAt
			}
		}

		if startsAt.IsZero() {
			startsAt = incidentReq.PublishedAt
		}
		if endsAt.IsZero() {
			endsAt = startsAt.Add(2 * time.Hour)
		}

		worstImpact := ""
		for _, summary := range incidentReq.StatusSummaries {
			worstImpact = getWorseStatus(worstImpact, summary.WorstComponentStatus)
		}

		maintenance := scheduledmaintenance.ScheduledMaintenance{
			Name:              incidentReq.Name,
			Link:              statusPageUrl + "/incidents/" + incidentReq.ID,
			StartsAt:          startsAt,
			EndsAt:            endsAt,
			ProviderImpact:    nullable.SetValue(worstImpact, worstImpact != ""),
			Impact:            nullable.SetValue(worstImpact, worstImpact != ""),
			ProviderID:        incidentReq.ID,
			ProviderCreatedAt: incidentReq.PublishedAt,
			Components:        []components.Component{},
		}

		for _, comp := range incidentReq.AffectedComponents {
			maintenance.Components = append(maintenance.Components, components.Component{
				ProviderID: comp.ComponentID,
			})
		}

		maintenance.Updates = make([]scheduledmaintenance.ScheduledMaintenanceUpdate, len(incidentReq.Updates))
		for idx, update := range incidentReq.Updates {
			maintenance.Updates[idx] = scheduledmaintenance.ScheduledMaintenanceUpdate{
				Description:                    update.MessageString,
				ScheduledMaintenanceProviderID: maintenance.ProviderID,
				ProviderID:                     update.ID,
				Status:                         statusnormalizer.NormalizeMaintenanceStatus(update.ToStatus),
				ProviderStatus:                 update.ToStatus,
				StatusTime:                     update.PublishedAt,
			}
		}

		maintenanceList = append(maintenanceList, maintenance)
	}

	return maintenanceList
}

func getWorseStatus(current, next string) string {
	severity := map[string]int{
		"none":                 0,
		"operational":          1,
		"degraded_performance": 2,
		"partial_outage":       3,
		"full_outage":          4,
	}

	if severity[next] > severity[current] {
		return next
	}
	return current
}
