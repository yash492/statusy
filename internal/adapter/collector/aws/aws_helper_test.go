package aws

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetchComponentsHelper(t *testing.T) {
	raw := `[
		{"service":"ec2-us-west-2","service_name":"Amazon EC2","region_id":"us-west-2","region_name":"Oregon"},
		{"service":"s3-us-west-2","service_name":"Amazon S3","region_id":"us-west-2","region_name":"Oregon"},
		{"service":"ec2-eu-west-1","service_name":"Amazon EC2","region_id":"eu-west-1","region_name":"Ireland"},
		{"service":"health","service_name":"AWS Health","region_name":null}
	]`

	var req ServicesReq
	require.NoError(t, json.Unmarshal([]byte(raw), &req))

	result := fetchComponentsHelper(req)
	assert.Empty(t, result.UngroupedComponents)
	require.Len(t, result.GroupedComponents, 3)

	// groups are sorted by provider id
	assert.Equal(t, "eu-west-1", result.GroupedComponents[0].ProviderID)
	assert.Equal(t, "eu-west-1 (Ireland)", result.GroupedComponents[0].Name)
	assert.Equal(t, "global", result.GroupedComponents[1].ProviderID)
	assert.Equal(t, "Global", result.GroupedComponents[1].Name)
	assert.Equal(t, "health", result.GroupedComponents[1].Components[0].ProviderID)

	oregon := result.GroupedComponents[2]
	assert.Equal(t, "us-west-2", oregon.ProviderID)
	require.Len(t, oregon.Components, 2)
	assert.Equal(t, "Amazon EC2", oregon.Components[0].Name)
	assert.Equal(t, "ec2-us-west-2", oregon.Components[0].ProviderID)
}

func TestFetchIncidentsHelper(t *testing.T) {
	raw := `{
		"ec2-us-west-2": [
			{
				"arn": "arn:aws:health:us-west-2::event/EC2/ISSUE/ABC",
				"summary": "[RESOLVED] Elevated Launch Failures",
				"status": "2",
				"date": "1794005465",
				"impacted_services": {
					"ec2-us-west-2": {"service_name": "Amazon EC2", "current": "0", "max": "2"},
					"batch-us-west-2": {"service_name": "AWS Batch", "current": "0", "max": "1"}
				},
				"event_log": [
					{"summary": "Elevated Launch Failures", "message": "Investigating.", "status": "1", "timestamp": 1794005465},
					{"summary": "Elevated Launch Failures", "message": "Hidden.", "status": "-1", "timestamp": 1794006000},
					{"summary": "[RESOLVED] Elevated Launch Failures", "message": "Resolved.", "status": "0", "timestamp": 1794009665}
				]
			}
		],
		"directconnect-eu-central-1": [
			{
				"arn": "arn:aws:health:eu-central-1::event/DX/ISSUE/DEF",
				"summary": "Connectivity Issues",
				"status": "3",
				"date": "1794000000",
				"event_log": [
					{"summary": "Connectivity Issues", "message": "Investigating.", "status": 3, "timestamp": 1794000000}
				]
			},
			{
				"arn": "arn:aws:health:eu-central-1::event/DX/ISSUE/NO_LOG",
				"summary": "No updates",
				"status": "1",
				"date": "1794000001",
				"event_log": []
			}
		],
		"s3-us-west-2": [
			{
				"arn": "arn:aws:health:us-west-2::event/EC2/ISSUE/ABC",
				"summary": "Duplicate under another key",
				"status": "1",
				"date": "1794005465",
				"event_log": [{"message": "dup", "status": "1", "timestamp": 1794005465}]
			}
		]
	}`

	var req HistoryEventsReq
	require.NoError(t, json.Unmarshal([]byte(raw), &req))

	result := fetchIncidentsHelper(req)
	// events without updates are skipped and the repeated arn is only returned once
	require.Len(t, result, 2)

	// keys are walked in sorted order
	dx := result[0]
	assert.Equal(t, "Connectivity Issues", dx.Name)
	assert.Equal(t, "arn:aws:health:eu-central-1::event/DX/ISSUE/DEF", dx.ProviderID)
	assert.Equal(t, "https://health.aws.amazon.com/health/status?eventID=arn%3Aaws%3Ahealth%3Aeu-central-1%3A%3Aevent%2FDX%2FISSUE%2FDEF", dx.Link)
	providerImpact, ok := dx.ProviderImpact.Get()
	assert.True(t, ok)
	assert.Equal(t, "3", providerImpact)
	impact, ok := dx.Impact.Get()
	assert.True(t, ok)
	assert.Equal(t, "disrupted", impact)
	// no impacted_services, falls back to the key the event is filed under
	require.Len(t, dx.Components, 1)
	assert.Equal(t, "directconnect-eu-central-1", dx.Components[0].ProviderID)
	// status given as a number instead of a string
	require.Len(t, dx.Updates, 1)
	assert.Equal(t, "investigating", dx.Updates[0].Status)
	assert.Equal(t, "3", dx.Updates[0].ProviderStatus)

	ec2 := result[1]
	assert.Equal(t, "Elevated Launch Failures", ec2.Name)
	assert.Equal(t, time.Unix(1794005465, 0).UTC(), ec2.ProviderCreatedAt)
	impact, ok = ec2.Impact.Get()
	assert.True(t, ok)
	assert.Equal(t, "degraded", impact)
	providerImpact, _ = ec2.ProviderImpact.Get()
	assert.Equal(t, "2", providerImpact)

	require.Len(t, ec2.Components, 2)
	assert.Equal(t, "batch-us-west-2", ec2.Components[0].ProviderID)
	assert.Equal(t, "AWS Batch", ec2.Components[0].Name)
	assert.Equal(t, "ec2-us-west-2", ec2.Components[1].ProviderID)

	// the hidden (-1) update is dropped
	require.Len(t, ec2.Updates, 2)
	assert.Equal(t, "investigating", ec2.Updates[0].Status)
	assert.Equal(t, "arn:aws:health:us-west-2::event/EC2/ISSUE/ABC#1794005465", ec2.Updates[0].ProviderID)
	assert.Equal(t, "resolved", ec2.Updates[1].Status)
	assert.Equal(t, "0", ec2.Updates[1].ProviderStatus)
	assert.Equal(t, "Resolved.", ec2.Updates[1].Description)
}

func TestFetchIncidentsHelperUnknownImpact(t *testing.T) {
	raw := `{"ec2-us-west-2": [{"arn": "a", "summary": "x", "status": "0", "date": "1",
		"event_log": [{"message": "m", "status": "0", "timestamp": 1}]}]}`

	var req HistoryEventsReq
	require.NoError(t, json.Unmarshal([]byte(raw), &req))

	result := fetchIncidentsHelper(req)
	require.Len(t, result, 1)
	_, ok := result[0].Impact.Get()
	assert.False(t, ok)
	_, ok = result[0].ProviderImpact.Get()
	assert.False(t, ok)
}
