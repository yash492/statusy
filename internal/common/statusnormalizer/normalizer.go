package statusnormalizer

import "strings"

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
