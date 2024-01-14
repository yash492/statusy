package squadcast

import (
	"strings"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const incidentName = "Flex Agents not receiving Legacy Chats"
const incidentDescription = "Our monitoring systems have detected a potential issue causing Flex Agents not receiving Legacy Chats. Our engineering team has been alerted and is actively investigating. We will update as soon as we have more information."
const incidentStatus = "Investigating"
const service = "Plivo"

var components = []string{"Phone", "SMS"}

type IncidentEvent struct {
	Message     string       `json:"message"`
	Description string       `json:"description"`
	Tags        IncidentTags `json:"tags"`
	Priority    string       `json:"priority"`
	Status      string       `json:"status"`
	EventID     string       `json:"event_id"`
}

type IncidentTags map[string]Tag

type Tag struct {
	Value string `json:"value"`
	Color string `json:"color,omitempty"`
}

func TestSquadcast(t *testing.T) {
	webhookURL := "https://api.squadcast.com/v2/incidents/api/dfc20b5ef2cd099d45779b7a1d908c30ce95053b"
	client := resty.New()

	event := IncidentEvent{
		Message:     "Plivo: " + incidentName,
		Description: incidentDescription,
		Tags: IncidentTags{
			"Service": Tag{
				Value: "Plivo",
			},
			"Components": Tag{
				Value: strings.Join(components, ", "),
				Color: "",
			},
			"Link": Tag{
				Value: "plivo.com",
			},
		},
		Status:  "trigger",
		EventID: uuid.New().String(),
	}
	_, err := client.R().SetBody(event).Post(webhookURL)
	assert.NoError(t, err)

}
