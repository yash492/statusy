package atlassianstatuspage

import (
	"encoding/json"
	"testing"

	"github.com/yash492/statusy/internal/domain/statuspage"
)

func TestCircleciComponents(t *testing.T) {

	type testAtlassianComponent struct {
		input struct {
			inputJson   string
			serviceSlug string
		}
		expectedOutput statuspage.AggregateComponents
	}

	tests := []testAtlassianComponent{
		{
			input: struct {
				inputJson   string
				serviceSlug string
			}{
				inputJson: `
				{
					"components": [
					{
						"id": "31vvcfzgyyzk",
						"name": "AWS",
						"status": "operational",
						"created_at": "2014-04-01T20:56:28.951Z",
						"updated_at": "2022-11-02T18:09:28.232Z",
						"position": 1,
						"description": "EC2 in us-east-1",
						"showcase": false,
						"start_date": null,
						"group_id": "5051gf6x40v1",
						"page_id": "6w4r0ttlx5ft",
						"group": false,
						"only_show_if_degraded": false
					},
					{
						"id": "8gmrsrb87lxr",
						"name": "Atlassian Bitbucket API",
						"status": "operational",
						"created_at": "2016-08-03T15:33:38.732Z",
						"updated_at": "2025-10-08T19:19:47.435Z",
						"position": 1,
						"description": null,
						"showcase": false,
						"start_date": null,
						"group_id": "s21fdjds0k15",
						"page_id": "6w4r0ttlx5ft",
						"group": false,
						"only_show_if_degraded": false
					},
					{
						"id": "9vxrymzc4lh5",
						"name": "Atlassian Bitbucket Source downloads",
						"status": "operational",
						"created_at": "2016-08-03T15:33:54.107Z",
						"updated_at": "2025-07-19T17:30:22.849Z",
						"position": 2,
						"description": null,
						"showcase": false,
						"start_date": null,
						"group_id": "s21fdjds0k15",
						"page_id": "6w4r0ttlx5ft",
						"group": false,
						"only_show_if_degraded": false
					},
					{
						"id": "sj405f2dg5ny",
						"name": "Google Cloud Platform Google Cloud DNS",
						"status": "operational",
						"created_at": "2016-09-01T00:28:16.059Z",
						"updated_at": "2023-12-05T01:49:58.175Z",
						"position": 2,
						"description": null,
						"showcase": false,
						"start_date": null,
						"group_id": "5051gf6x40v1",
						"page_id": "6w4r0ttlx5ft",
						"group": false,
						"only_show_if_degraded": false
					},
					{
						"id": "k7z3xkf61sff",
						"name": "CircleCI Releases",
						"status": "operational",
						"created_at": "2024-02-19T17:40:30.067Z",
						"updated_at": "2025-05-02T13:48:42.246Z",
						"position": 12,
						"description": "Subsystem responsible for tracking and managing releases.",
						"showcase": true,
						"start_date": "2024-02-19",
						"group_id": null,
						"page_id": "6w4r0ttlx5ft",
						"group": false,
						"only_show_if_degraded": false
					},

					{
						"id": "gx397f9wvq4w",
						"name": "Billing \u0026 Account",
						"status": "operational",
						"created_at": "2020-05-22T18:12:56.408Z",
						"updated_at": "2025-07-12T20:00:21.160Z",
						"position": 14,
						"description": "Payment processing and billing account changes.",
						"showcase": false,
						"start_date": null,
						"group_id": null,
						"page_id": "6w4r0ttlx5ft",
						"group": false,
						"only_show_if_degraded": false
					},
					{
						"id": "5051gf6x40v1",
						"name": "CircleCI Dependencies",
						"status": "operational",
						"created_at": "2020-10-23T23:02:58.992Z",
						"updated_at": "2025-05-02T13:48:42.302Z",
						"position": 15,
						"description": "Third party services that CircleCI infrastructure depends on.",
						"showcase": false,
						"start_date": null,
						"group_id": null,
						"page_id": "6w4r0ttlx5ft",
						"group": true,
						"only_show_if_degraded": false,
						"components": [
							"31vvcfzgyyzk",
							"sj405f2dg5ny"
						]
					},
					{
						"id": "s21fdjds0k15",
						"name": "Upstream Services",
						"status": "operational",
						"created_at": "2020-10-23T23:23:34.529Z",
						"updated_at": "2025-05-02T13:48:42.342Z",
						"position": 16,
						"description": "Third party services that may impact jobs.",
						"showcase": false,
						"start_date": null,
						"group_id": null,
						"page_id": "6w4r0ttlx5ft",
						"group": true,
						"only_show_if_degraded": false,
						"components": [
							"8gmrsrb87lxr",
							"9vxrymzc4lh5"
						]
					}
				]
			}	`,
				serviceSlug: circleciSlug,
			},
			expectedOutput: statuspage.AggregateComponents{
				GroupedComponents: []statuspage.ComponentGroup{
					{
						Name:       "CircleCI Dependencies",
						ProviderID: "5051gf6x40v1",
						Components: []statuspage.Component{
							{Name: "AWS", ProviderID: "31vvcfzgyyzk"},
							{Name: "Google Cloud Platform Google Cloud DNS", ProviderID: "sj405f2dg5ny"},
						},
					},
					{
						Name:       "Upstream Services",
						ProviderID: "s21fdjds0k15",
						Components: []statuspage.Component{
							{Name: "Atlassian Bitbucket API", ProviderID: "8gmrsrb87lxr"},
							{Name: "Atlassian Bitbucket Source downloads", ProviderID: "9vxrymzc4lh5"},
						},
					},
				},
				UngroupedComponents: []statuspage.Component{
					{Name: "CircleCI Releases", ProviderID: "k7z3xkf61sff"},
					{Name: "Billing & Account", ProviderID: "gx397f9wvq4w"},
				},
			},
		},
	}

	for _, tt := range tests {
		var atlassianComponents atlassianComponentsReq
		err := json.Unmarshal([]byte(tt.input.inputJson), &atlassianComponents)
		if err != nil {
			t.Fatalf("unmarshalling circleci components went wrong %s", err.Error())
		}

		components := fetchComponentsHelper(atlassianComponents, tt.input.serviceSlug)
		resultOutputComponentMap := map[string]statuspage.Component{}
		for _, ungroupedComponent := range components.UngroupedComponents {
			resultOutputComponentMap[ungroupedComponent.ProviderID] = ungroupedComponent
		}

		for _, expectedUngroupedComponent := range tt.expectedOutput.UngroupedComponents {
			resultOutputComponent, ok := resultOutputComponentMap[expectedUngroupedComponent.ProviderID]
			if !ok {
				t.Fatalf("component %s (provider id: %s) does not exists in expected components", expectedUngroupedComponent.Name, expectedUngroupedComponent.ProviderID)
			}

			if resultOutputComponent.Name != expectedUngroupedComponent.Name {
				t.Fatalf("expected component name is %s, got %s", expectedUngroupedComponent.Name, resultOutputComponent.Name)
			}

			if resultOutputComponent.ProviderID != expectedUngroupedComponent.ProviderID {
				t.Fatalf("expected component provider_id is %s, got %s", expectedUngroupedComponent.ProviderID, resultOutputComponent.ProviderID)
			}

		}

		resultOutputComponentMapForComponentGroup := map[string]map[string]statuspage.Component{}
		resultOutputComponentGroupMap := map[string]statuspage.ComponentGroup{}

		for _, groupedComponent := range components.GroupedComponents {
			resultOutputComponentGroupMap[groupedComponent.ProviderID] = groupedComponent
			for _, component := range groupedComponent.Components {
				output, ok := resultOutputComponentMapForComponentGroup[groupedComponent.ProviderID]
				if !ok {
					resultOutputComponentMapForComponentGroup[groupedComponent.ProviderID] = map[string]statuspage.Component{
						component.ProviderID: component,
					}
				} else {
					output[component.ProviderID] = component
				}
			}
		}

		for _, component := range tt.expectedOutput.GroupedComponents {
			resultOutputComponentGroup, ok := resultOutputComponentGroupMap[component.ProviderID]
			if !ok {
				t.Fatalf("component group %s (provider id: %s) does not exists in expected components", resultOutputComponentGroup.Name, resultOutputComponentGroup.ProviderID)
			}

			if resultOutputComponentGroup.Name != component.Name {
				t.Fatalf("expected component name is %s, got %s", component.Name, resultOutputComponentGroup.Name)

			}

			if resultOutputComponentGroup.ProviderID != component.ProviderID {
				t.Fatalf("expected component provider_id is %s, got %s", component.ProviderID, resultOutputComponentGroup.ProviderID)
			}

			for _, expectedChildComponent := range component.Components {
				resultChildComponents := resultOutputComponentMapForComponentGroup[component.ProviderID]
				resultOutputComponent, ok := resultChildComponents[expectedChildComponent.ProviderID]
				if !ok {
					t.Fatalf("component %s (provider id: %s) does not exists in expected components", expectedChildComponent.Name, expectedChildComponent.ProviderID)
				}

				if resultOutputComponent.Name != expectedChildComponent.Name {
					t.Fatalf("expected component name is %s, got %s", expectedChildComponent.Name, resultOutputComponent.Name)
				}

				if resultOutputComponent.ProviderID != expectedChildComponent.ProviderID {
					t.Fatalf("expected component provider_id is %s, got %s", expectedChildComponent.ProviderID, resultOutputComponent.ProviderID)
				}

			}

		}
	}

}
