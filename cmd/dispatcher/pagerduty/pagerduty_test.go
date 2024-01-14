package pagerduty

import (
	"context"
	"strings"
	"testing"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/stretchr/testify/assert"
	"github.com/yash492/statusy/pkg/types"
)

const incidentName = "Flex Agents not receiving Legacy Chats"
const incidentDescription = "Our monitoring systems have detected a potential issue causing Flex Agents not receiving Legacy Chats. Our engineering team has been alerted and is actively investigating. We will update as soon as we have more information."
const incidentStatus = "Investigating"
const service = "Plivo"

var components = []string{"Phone", "SMS"}

func TestPagerduty(t *testing.T) {
	ctx := context.Background()
	a := assert.New(t)
	routingKey := "fa3015e76a8f4c0ed0f8af05d1d999e5"
	resp, err := pagerduty.ManageEventWithContext(ctx, pagerduty.V2Event{
		RoutingKey: routingKey,
		Action:     "trigger",
		Client:     "Plivo",
		ClientURL:  "https://status.plivo.com",
		Payload: &pagerduty.V2Payload{
			Summary:  incidentName,
			Source:   "Statusy",
			Severity: "warning",

			Component: strings.Join(components, ", "),
			Details: types.JSON{
				"More Info": incidentDescription,
			},
		},
	})

	t.Log(resp, "\n")

	a.NoError(err)

}
