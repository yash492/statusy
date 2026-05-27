package incidentio_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yash492/statusy/internal/adapter/collector/incidentio"
)

func TestFetchComponentsHelper(t *testing.T) {
	// Read the sample status JSON
	path := filepath.Join("../../../../", "tmp", "openai_status.json")
	bytes, err := os.ReadFile(path)
	if err != nil {
		t.Skip("Skipping test: tmp/openai_status.json not found")
		return
	}

	var req incidentio.StatusReq
	err = json.Unmarshal(bytes, &req)
	assert.NoError(t, err)

	componentsResult := incidentio.FetchComponentsHelper(req)
	assert.NotEmpty(t, componentsResult.GroupedComponents)

	// Verify group APIs exists
	var foundAPIs bool
	for _, g := range componentsResult.GroupedComponents {
		if g.Name == "APIs" {
			foundAPIs = true
			assert.NotEmpty(t, g.Components)
			assert.Equal(t, "01K5H8S53SY1KMS4GQMNMQM1K5", g.ProviderID)
		}
	}
	assert.True(t, foundAPIs)
}

func TestFetchIncidentsHelper(t *testing.T) {
	// Read the sample incidents JSON
	path := filepath.Join("../../../../", "tmp", "openai_incidents.json")
	bytes, err := os.ReadFile(path)
	if err != nil {
		t.Skip("Skipping test: tmp/openai_incidents.json not found")
		return
	}

	var req incidentio.IncidentsReq
	err = json.Unmarshal(bytes, &req)
	assert.NoError(t, err)

	incidentsResult := incidentio.FetchIncidentsHelper(req, "https://status.openai.com")
	assert.NotEmpty(t, incidentsResult)

	for _, incident := range incidentsResult {
		assert.NotEmpty(t, incident.Name)
		assert.NotEmpty(t, incident.Link)
		assert.NotEmpty(t, incident.ProviderID)
		assert.False(t, incident.ProviderCreatedAt.IsZero())
		for _, update := range incident.Updates {
			assert.NotEmpty(t, update.Description)
			assert.NotEmpty(t, update.ProviderID)
			assert.NotEmpty(t, update.Status)
			assert.False(t, update.StatusTime.IsZero())
		}
	}
}

func TestFetchIncidentsAndMaintenances(t *testing.T) {
	mockJSON := `{
  "incidents": [
    {
      "updates": [
        {
          "published_at": "2026-01-30T10:42:56.002Z",
          "id": "01KG77ZCW3DBEXXB00VHC5GVZV",
          "message_string": "Some customers are unable to access the Invoices page. We have identified the issue and are working on a fix.",
          "to_status": "identified",
          "component_statuses": [
            {
              "component_id": "01JEV56RY027PA4NYX5H2RNAMC",
              "status": "degraded_performance"
            }
          ],
          "automated_update": false
        }
      ],
      "component_impacts": [
        {
          "start_at": "2026-01-30T10:42:56.002Z",
          "end_at": "2026-01-30T12:38:56.689Z",
          "id": "01KG77ZCW3QHCS75F58Y381918",
          "component_id": "01JEV56RY027PA4NYX5H2RNAMC",
          "status_page_incident_id": "01KG77ZCW3ZCZBXCRP47BT0KFG",
          "status": "degraded_performance"
        }
      ],
      "status_summaries": [
        {
          "start_at": "2026-01-30T10:42:56.002Z",
          "end_at": "2026-01-30T11:21:05.2Z",
          "worst_component_status": "degraded_performance"
        }
      ],
      "published_at": "2026-01-30T10:42:56.002Z",
      "id": "01KG77ZCW3ZCZBXCRP47BT0KFG",
      "status_page_id": "01JEV56RY0P76Q6QPHRZN5SJ9D",
      "name": "Invoices page not loading for some customers",
      "status": "resolved",
      "affected_components": [
        {
          "component_id": "01JEV56RY027PA4NYX5H2RNAMC",
          "status": "degraded_performance",
          "current_status": "operational"
        }
      ],
      "type": "incident"
    },
    {
      "updates": [
        {
          "published_at": "2026-01-21T14:43:41.913Z",
          "id": "01KFGG5S6T1Q0YQFDB2YWSFBAD",
          "message_string": "We're reaching out to inform you of a scheduled maintenance window...",
          "to_status": "maintenance_scheduled",
          "component_statuses": [
            {
              "component_id": "01JEV56RY04T2614381RR4X5AJ",
              "status": "under_maintenance"
            }
          ],
          "automated_update": false
        }
      ],
      "component_impacts": [
        {
          "start_at": "2026-01-25T00:00:00Z",
          "end_at": "2026-01-25T02:00:00Z",
          "id": "01KFS9X4QSQXBEZYZ8M4PCKVJA",
          "component_id": "01JEV56RY0GV18JSH8N9GECJ1G",
          "status_page_incident_id": "01KFGG5S6TRSJJNNBB3TQATY24",
          "status": "under_maintenance"
        }
      ],
      "status_summaries": [
        {
          "start_at": "2026-01-25T00:00:00Z",
          "end_at": "2026-01-25T02:00:00Z",
          "worst_component_status": "under_maintenance"
        }
      ],
      "published_at": "2026-01-21T14:43:41.913Z",
      "id": "01KFGG5S6TRSJJNNBB3TQATY24",
      "status_page_id": "01JEV56RY0P76Q6QPHRZN5SJ9D",
      "name": "Scheduled Pleo Maintenance",
      "status": "maintenance_complete",
      "affected_components": [
        {
          "component_id": "01JEV56RY0GV18JSH8N9GECJ1G",
          "status": "under_maintenance",
          "current_status": "operational"
        }
      ],
      "type": "maintenance"
    }
  ]
}`

	var req incidentio.IncidentsReq
	err := json.Unmarshal([]byte(mockJSON), &req)
	assert.NoError(t, err)

	incidentsResult := incidentio.FetchIncidentsHelper(req, "https://status.pleo.io")
	assert.Len(t, incidentsResult, 1)
	assert.Equal(t, "01KG77ZCW3ZCZBXCRP47BT0KFG", incidentsResult[0].ProviderID)
	assert.Equal(t, "Invoices page not loading for some customers", incidentsResult[0].Name)

	maintenancesResult := incidentio.FetchScheduledMaintenancesHelper(req, "https://status.pleo.io")
	assert.Len(t, maintenancesResult, 1)
	assert.Equal(t, "01KFGG5S6TRSJJNNBB3TQATY24", maintenancesResult[0].ProviderID)
	assert.Equal(t, "Scheduled Pleo Maintenance", maintenancesResult[0].Name)
}
