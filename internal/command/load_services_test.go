package command

import (
	"context"

	"github.com/goccy/go-yaml"
	postgres "github.com/yash492/statusy/internal/adapter/postgres/services"
	domainservices "github.com/yash492/statusy/internal/domain/services"
)

func (t *TestSuite) TestLoadServices() {

	tests := []struct {
		inputYaml      []map[string]any
		expectedResult []domainservices.ServiceResult
	}{
		{
			inputYaml: []map[string]any{
				{
					"name":                      "Plivo",
					"slug":                      "plivo",
					"incidents_url":             "https://status.plivo.com/api/v2/incidents.json",
					"components_url":            "https://status.plivo.com/api/v2/components.json",
					"schedule_maintenances_url": "https://status.plivo.com/api/v2/scheduled-maintenances.json",
					"provider_type":             "atlassian",
				},
				{
					"name":                      "Circle CI",
					"slug":                      "circleci",
					"incidents_url":             "https://status.circleci.com/api/v2/incidents.json",
					"components_url":            "https://status.circleci.com/api/v2/components.json",
					"schedule_maintenances_url": "https://status.circleci.com/api/v2/scheduled-maintenances.json",
					"provider_type":             "atlassian",
				},
			},
			expectedResult: []domainservices.ServiceResult{
				{
					Name:                    "Plivo",
					Slug:                    "plivo",
					IncidentsUrl:            "https://status.circleci.com/api/v2/incidents.json",
					ScheduleMaintenancesUrl: "https://status.circleci.com/api/v2/scheduled-maintenances.json",
					ComponentsUrl:           "https://status.circleci.com/api/v2/components.json",
					ProviderType:            "atlassian",
				},
				{
					Name:                    "Circle CI",
					Slug:                    "circleci",
					IncidentsUrl:            "https://status.circleci.com/api/v2/incidents.json",
					ScheduleMaintenancesUrl: "https://status.circleci.com/api/v2/scheduled-maintenances.json",
					ComponentsUrl:           "https://status.circleci.com/api/v2/components.json",
					ProviderType:            "atlassian",
				},
			},
		},
	}

	repo := postgres.NewPostgresServiceRepository(t.Logger, t.TestDb, t.TestDb)

	for _, tt := range tests {

		bytes, err := yaml.Marshal(tt.inputYaml)
		if err != nil {
			t.T().Fatalf("could not marshal services map %s", err)
		}
		cmd := NewLoadServiceCmd(t.Logger, bytes, repo)
		services, err := cmd.Execute(context.Background())

		if err != nil {
			t.T().Fatalf("unable to insert and retrieve services. err: %s", err.Error())
		}

		if len(services) != len(tt.expectedResult) {
			t.T().Fatalf("mismatch between result len and expected service len. got=%d expected=%d", len(services), len(tt.expectedResult))
		}

		// for _, service := range services {
		// 	if service.ID == 0 {
		// 		t.T().Fatal("service ID cannot be zero")
		// 	}

		// }
	}

}
