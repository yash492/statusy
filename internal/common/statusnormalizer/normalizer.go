package statusnormalizer

import "strings"

type NormalizedStatus interface {
	IsResolved() bool
	IsInitial() bool
	String() string
	ForDisplay() string
}

type IncidentStatus string

const (
	IncidentStatusInvestigating IncidentStatus = "investigating"
	IncidentStatusIdentified    IncidentStatus = "identified"
	IncidentStatusMonitoring    IncidentStatus = "monitoring"
	IncidentStatusResolved      IncidentStatus = "resolved"
)

func (s IncidentStatus) IsResolved() bool { return s == IncidentStatusResolved }
func (s IncidentStatus) IsInitial() bool  { return s == IncidentStatusInvestigating }
func (s IncidentStatus) String() string   { return string(s) }
func (s IncidentStatus) ForDisplay() string {
	switch s {
	case IncidentStatusInvestigating:
		return "Investigating"
	case IncidentStatusIdentified:
		return "Identified"
	case IncidentStatusMonitoring:
		return "Monitoring"
	case IncidentStatusResolved:
		return "Resolved"
	default:
		return string(s)
	}
}

func ParseIncidentStatus(s string) IncidentStatus { return IncidentStatus(s) }

type MaintenanceStatus string

const (
	MaintenanceStatusScheduled  MaintenanceStatus = "scheduled"
	MaintenanceStatusInProgress MaintenanceStatus = "in_progress"
	MaintenanceStatusVerifying  MaintenanceStatus = "verifying"
	MaintenanceStatusCompleted  MaintenanceStatus = "completed"
)

func (s MaintenanceStatus) IsResolved() bool { return s == MaintenanceStatusCompleted }
func (s MaintenanceStatus) IsInitial() bool  { return s == MaintenanceStatusScheduled }
func (s MaintenanceStatus) String() string   { return string(s) }
func (s MaintenanceStatus) ForDisplay() string {
	switch s {
	case MaintenanceStatusScheduled:
		return "Scheduled"
	case MaintenanceStatusInProgress:
		return "In Progress"
	case MaintenanceStatusVerifying:
		return "Verifying"
	case MaintenanceStatusCompleted:
		return "Completed"
	default:
		return string(s)
	}
}

func ParseMaintenanceStatus(s string) MaintenanceStatus { return MaintenanceStatus(s) }

// NormalizeIncidentStatus standardizes incident states to: investigating, identified, monitoring, resolved
func NormalizeIncidentStatus(status string) string {
	s := strings.ToLower(strings.TrimSpace(status))
	switch s {
	case "investigating":
		return "investigating"
	case "identified", "identified_incident":
		return "identified"
	case "monitoring", "monitoring_incident":
		return "monitoring"
	case "resolved", "completed", "postmortem":
		return "resolved"
	default:
		return "investigating"
	}
}

// NormalizeMaintenanceStatus standardizes scheduled maintenance states to: scheduled, in_progress, verifying, completed
func NormalizeMaintenanceStatus(status string) string {
	s := strings.ToLower(strings.TrimSpace(status))
	switch s {
	case "scheduled", "planned", "maintenance_scheduled":
		return "scheduled"
	case "in_progress", "active", "maintenance_in_progress":
		return "in_progress"
	case "verifying", "monitoring":
		return "verifying"
	case "completed", "maintenance_complete", "postmortem":
		return "completed"
	default:
		return "scheduled"
	}
}
