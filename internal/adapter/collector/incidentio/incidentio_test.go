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
