package statusnormalizer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeIncidentStatus(t *testing.T) {
	assert.Equal(t, "investigating", NormalizeIncidentStatus("INVESTIGATING"))
	assert.Equal(t, "identified", NormalizeIncidentStatus("identified_incident"))
	assert.Equal(t, "monitoring", NormalizeIncidentStatus("monitoring"))
	assert.Equal(t, "resolved", NormalizeIncidentStatus("completed"))
	assert.Equal(t, "resolved", NormalizeIncidentStatus("postmortem"))
	assert.Equal(t, "investigating", NormalizeIncidentStatus("unknown_status"))
}

func TestNormalizeMaintenanceStatus(t *testing.T) {
	assert.Equal(t, "scheduled", NormalizeMaintenanceStatus("PLANNED"))
	assert.Equal(t, "scheduled", NormalizeMaintenanceStatus("maintenance_scheduled"))
	assert.Equal(t, "in_progress", NormalizeMaintenanceStatus("active"))
	assert.Equal(t, "in_progress", NormalizeMaintenanceStatus("maintenance_in_progress"))
	assert.Equal(t, "verifying", NormalizeMaintenanceStatus("monitoring"))
	assert.Equal(t, "completed", NormalizeMaintenanceStatus("completed"))
	assert.Equal(t, "completed", NormalizeMaintenanceStatus("maintenance_complete"))
	assert.Equal(t, "completed", NormalizeMaintenanceStatus("postmortem"))
	assert.Equal(t, "scheduled", NormalizeMaintenanceStatus("unknown_status"))
}
