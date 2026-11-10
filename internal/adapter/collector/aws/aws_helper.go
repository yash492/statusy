package aws

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/yash492/statusy/internal/common/nullable"
	"github.com/yash492/statusy/internal/common/statusnormalizer"
	"github.com/yash492/statusy/internal/domain/components"
	"github.com/yash492/statusy/internal/domain/incidents"
)

const (
	globalGroupID   = "global"
	globalGroupName = "Global"
	resolvedPrefix  = "[RESOLVED]"
	// providerStatusResolved is the event_log status AWS uses once an event is over.
	providerStatusResolved = "0"
)

// awsImpacts maps the AWS event status to the label used by the AWS Health Dashboard.
var awsImpacts = map[string]string{
	"1": "impacted",
	"2": "degraded",
	"3": "disrupted",
}

// ServicesReq represents services.json: one row per service and region.
type ServicesReq []ServiceRow

type ServiceRow struct {
	Service     string  `json:"service"`
	ServiceName string  `json:"service_name"`
	RegionID    string  `json:"region_id"`
	RegionName  *string `json:"region_name"`
}

// HistoryEventsReq represents historyevents.json, keyed by "<service>-<region>".
// Active events are part of it as well, so currentevents is not needed.
type HistoryEventsReq map[string][]HistoryEvent

type HistoryEvent struct {
	ARN              string                     `json:"arn"`
	Summary          string                     `json:"summary"`
	Status           json.Number                `json:"status"`
	Date             json.Number                `json:"date"`
	EventLog         []EventLog                 `json:"event_log"`
	ImpactedServices map[string]ImpactedService `json:"impacted_services"`
}

type EventLog struct {
	Message   string      `json:"message"`
	Status    json.Number `json:"status"`
	Timestamp int64       `json:"timestamp"`
}

type ImpactedService struct {
	ServiceName string `json:"service_name"`
}

func fetchComponentsHelper(req ServicesReq) components.AggregateComponents {
	groups := map[string]*components.ComponentGroup{}

	for _, row := range req {
		groupID, groupName := globalGroupID, globalGroupName
		if row.RegionID != "" {
			groupID = row.RegionID
			groupName = fmt.Sprintf("%s (%s)", row.RegionID, lo.FromPtr(row.RegionName))
		}

		group, ok := groups[groupID]
		if !ok {
			group = &components.ComponentGroup{Name: groupName, ProviderID: groupID}
			groups[groupID] = group
		}
		group.Components = append(group.Components, components.Component{
			Name:       row.ServiceName,
			ProviderID: row.Service,
		})
	}

	groupedComponents := lo.MapToSlice(groups, func(_ string, group *components.ComponentGroup) components.ComponentGroup {
		sort.Slice(group.Components, func(i, j int) bool {
			return group.Components[i].Name < group.Components[j].Name
		})
		return *group
	})
	sort.Slice(groupedComponents, func(i, j int) bool {
		return groupedComponents[i].ProviderID < groupedComponents[j].ProviderID
	})

	return components.AggregateComponents{
		UngroupedComponents: []components.Component{},
		GroupedComponents:   groupedComponents,
	}
}

func fetchIncidentsHelper(req HistoryEventsReq) []incidents.Incident {
	incidentList := []incidents.Incident{}
	seen := map[string]struct{}{}

	// Map iteration order is random, so walk the keys sorted to keep the output stable.
	keys := lo.Keys(req)
	sort.Strings(keys)

	for _, key := range keys {
		for _, event := range req[key] {
			// The dashboard skips events without updates and hidden (-1) updates.
			logs := lo.Filter(event.EventLog, func(l EventLog, _ int) bool {
				return l.Status.String() != "-1"
			})
			if len(logs) == 0 {
				continue
			}
			if _, ok := seen[event.ARN]; ok {
				continue
			}
			seen[event.ARN] = struct{}{}

			createdAt, _ := event.Date.Int64()
			impact, hasImpact := awsImpacts[event.Status.String()]

			incident := incidents.Incident{
				Name:              strings.TrimSpace(strings.TrimPrefix(event.Summary, resolvedPrefix)),
				Link:              fmt.Sprintf("%s?eventID=%s", statusPageUrl, url.QueryEscape(event.ARN)),
				ProviderImpact:    nullable.SetValue(event.Status.String(), hasImpact),
				Impact:            nullable.SetValue(impact, hasImpact),
				ProviderID:        event.ARN,
				ProviderCreatedAt: time.Unix(createdAt, 0).UTC(),
			}

			incident.Updates = lo.Map(logs, func(l EventLog, _ int) incidents.IncidentUpdate {
				return incidents.IncidentUpdate{
					Description:        l.Message,
					IncidentProviderID: incident.ProviderID,
					ProviderID:         fmt.Sprintf("%s#%d", event.ARN, l.Timestamp),
					Status:             normalizeEventLogStatus(l.Status.String()),
					ProviderStatus:     l.Status.String(),
					StatusTime:         time.Unix(l.Timestamp, 0).UTC(),
				}
			})

			incident.Components = incidentComponents(key, event)
			incidentList = append(incidentList, incident)
		}
	}

	return incidentList
}

// normalizeEventLogStatus converts the AWS severity-style status into a lifecycle status.
// AWS does not publish identified/monitoring states, so everything unresolved is investigating.
func normalizeEventLogStatus(status string) string {
	if status == providerStatusResolved {
		return statusnormalizer.NormalizeIncidentStatus("resolved")
	}
	return statusnormalizer.NormalizeIncidentStatus("investigating")
}

// incidentComponents uses the impacted services of the event. Some events do not list them,
// in that case the "<service>-<region>" key the event is filed under is used.
func incidentComponents(key string, event HistoryEvent) []components.Component {
	if len(event.ImpactedServices) == 0 {
		return []components.Component{{Name: key, ProviderID: key}}
	}

	impactedKeys := lo.Keys(event.ImpactedServices)
	sort.Strings(impactedKeys)

	return lo.Map(impactedKeys, func(impactedKey string, _ int) components.Component {
		return components.Component{
			Name:       event.ImpactedServices[impactedKey].ServiceName,
			ProviderID: impactedKey,
		}
	})
}
