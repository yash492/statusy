package notification

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yash492/statusy/internal/domain/notifications"
)

type LiveTestConfigs struct {
	SlackWebhookURL      string
	DiscordWebhookURL    string
	MsTeamsWebhookURL    string
	PagerDutyRoutingKey  string
	SolarwindsWebhookURL string
	WebhookURL           string
}

// Set these fields to test live webhook dispatching directly
var liveConfigs = LiveTestConfigs{
	SlackWebhookURL:      "https://hooks.slack.com/services/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",
	DiscordWebhookURL:    "https://discord.com/api/webhooks/1513573666166935746/10pbGuOCr701_in9bDHd8vBujs3DmLmgKAslKi6jcPwvQku8XMVZI-R2AgDVzPO8vnje",
	MsTeamsWebhookURL:    "",
	PagerDutyRoutingKey:  "",
	SolarwindsWebhookURL: "",
	WebhookURL:           "",
}

func TestNotifierChannelsLive(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	notifier := NewHttpNotifier(logger)

	incidentDetails := notifications.IncidentNotificationDetails{
		IncidentID:  12,
		UpdateID:    34,
		Title:       "Major Outage",
		Status:      "Investigating",
		Description: "Database connection failed",
		ServiceName: "Auth Service",
		Components: []notifications.NotificationComponent{
			{Name: "DB Engine"},
		},
		UpdatedAt: time.Now(),
		Link:      "https://status.example.com/incidents/12",
	}

	maintenanceDetails := notifications.MaintenanceNotificationDetails{
		MaintenanceID: 56,
		UpdateID:      78,
		Title:         "Database Upgrade",
		Status:        "scheduled",
		Description:   "Upgrading DB to v16",
		ServiceName:   "DB Cluster",
		Components: []notifications.NotificationComponent{
			{Name: "DB Master"},
		},
		StartTime: time.Now(),
		EndTime:   time.Now().Add(2 * time.Hour),
		UpdatedAt: time.Now(),
		Link:      "https://status.example.com/maintenance/56",
	}

	t.Run("Slack Webhook Dispatch", func(t *testing.T) {
		if liveConfigs.SlackWebhookURL == "" {
			t.Skip("Slack URL not set")
		}

		ch := notifications.ViewNotification{
			Type:   notifications.NotificationTypeSlack,
			Config: json.RawMessage(`{"webhook_url": "` + liveConfigs.SlackWebhookURL + `"}`),
		}

		extID, err := notifier.SendIncident(context.Background(), ch, true, incidentDetails, "")
		assert.NoError(t, err)
		assert.Empty(t, extID)
	})

	t.Run("Discord Webhook Dispatch - POST and PATCH", func(t *testing.T) {
		if liveConfigs.DiscordWebhookURL == "" {
			t.Skip("Discord URL not set")
		}

		ch := notifications.ViewNotification{
			Type:   notifications.NotificationTypeDiscord,
			Config: json.RawMessage(`{"webhook_url": "` + liveConfigs.DiscordWebhookURL + `"}`),
		}

		// 1. Initial message (POST)
		extID, err := notifier.SendIncident(context.Background(), ch, true, incidentDetails, "")
		assert.NoError(t, err)
		assert.NotEmpty(t, extID)

		// 2. Subsequent update (PATCH)
		extID2, err := notifier.SendIncident(context.Background(), ch, false, incidentDetails, extID)
		assert.NoError(t, err)
		assert.Equal(t, extID, extID2)
	})

	t.Run("MS Teams Webhook Dispatch", func(t *testing.T) {
		if liveConfigs.MsTeamsWebhookURL == "" {
			t.Skip("MS Teams URL not set")
		}

		ch := notifications.ViewNotification{
			Type:   notifications.NotificationTypeMsTeams,
			Config: json.RawMessage(`{"webhook_url": "` + liveConfigs.MsTeamsWebhookURL + `"}`),
		}

		extID, err := notifier.SendIncident(context.Background(), ch, true, incidentDetails, "")
		assert.NoError(t, err)
		assert.Empty(t, extID)
	})

	t.Run("PagerDuty V2 Events Dispatch", func(t *testing.T) {
		if liveConfigs.PagerDutyRoutingKey == "" {
			t.Skip("PagerDuty Routing Key not set")
		}

		ch := notifications.ViewNotification{
			Type:   notifications.NotificationTypePagerDuty,
			Config: json.RawMessage(`{"routing_key": "` + liveConfigs.PagerDutyRoutingKey + `"}`),
		}

		extID, err := notifier.SendIncident(context.Background(), ch, true, incidentDetails, "")
		assert.NoError(t, err)
		assert.NotEmpty(t, extID)
	})

	t.Run("SolarWinds Incident Response Webhook", func(t *testing.T) {
		if liveConfigs.SolarwindsWebhookURL == "" {
			t.Skip("Solarwinds URL not set")
		}

		ch := notifications.ViewNotification{
			Type:   notifications.NotificationTypeSolarwindsIncidentResponse,
			Config: json.RawMessage(`{"webhook_url": "` + liveConfigs.SolarwindsWebhookURL + `"}`),
		}

		extID, err := notifier.SendIncident(context.Background(), ch, true, incidentDetails, "")
		assert.NoError(t, err)
		assert.Empty(t, extID)
	})

	t.Run("Custom Webhook Dispatch", func(t *testing.T) {
		if liveConfigs.WebhookURL == "" {
			t.Skip("Webhook URL not set")
		}

		ch := notifications.ViewNotification{
			Type:   notifications.NotificationTypeWebhook,
			Config: json.RawMessage(`{"url": "` + liveConfigs.WebhookURL + `", "headers": {"X-Test-Header": "hello-world"}}`),
		}

		extID, err := notifier.SendMaintenance(context.Background(), ch, true, maintenanceDetails, "")
		assert.NoError(t, err)
		assert.Empty(t, extID)
	})
}
