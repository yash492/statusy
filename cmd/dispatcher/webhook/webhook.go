package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/jackc/pgx"
	"github.com/samber/lo"
	"github.com/yash492/statusy/pkg/domain"
	"github.com/yash492/statusy/pkg/types"
)

type IncidentOpenWorker struct{}
type IncidentInProgressWorker struct{}
type IncidentClosedWorker struct{}

type webhookResponse struct {
	ServiceID                    uint               `json:"service_id"`
	ServiceName                  string             `json:"service_name"`
	IncidentID                   uint               `json:"incident_id"`
	IncidentName                 string             `json:"incident_name"`
	IncidentLink                 string             `json:"incident_link"`
	IncidentImpact               string             `json:"incident_impact"`
	IncidentUpdate               string             `json:"incident_update"`
	IncidentUpdateID             uint               `json:"incident_update_id"`
	IncidentUpdateProviderStatus string             `json:"incident_update_provider_status"`
	IncidentUpdateStatus         string             `json:"incident_update_status"`
	IncidentUpdateComponents     []webhookComponent `json:"incident_update_components"`
	EventType                    string             `json:"event_type"`
	IncidentUpdateStatusTime     time.Time          `json:"incident_update_status_time"`
}

type webhookComponent struct {
	Name string `json:"name"`
	ID   uint   `json:"id"`
}

func dispatchWebhookMsg(event types.WorkerEvent) error {
	webhook, err := domain.WebhookExtension.Get()
	if err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return nil
		}
		return err
	}

	components := lo.Map(event.Components, func(component types.ComponentsForWorker, _ int) webhookComponent {
		return webhookComponent{
			Name: component.Name,
			ID:   component.ID,
		}
	})

	payload := webhookResponse{
		ServiceID:                    event.ServiceID,
		ServiceName:                  event.ServiceName,
		IncidentID:                   event.IncidentID,
		IncidentName:                 event.IncidentName,
		IncidentLink:                 event.IncidentLink,
		IncidentImpact:               event.IncidentImpact,
		IncidentUpdate:               event.IncidentUpdate,
		IncidentUpdateID:             event.IncidentUpdateID,
		IncidentUpdateProviderStatus: event.IncidentUpdateProviderStatus,
		IncidentUpdateStatus:         event.IncidentUpdateStatus,
		IncidentUpdateComponents:     components,
		EventType:                    event.EventType,
		IncidentUpdateStatusTime:     event.IncidentUpdateStatusTime,
	}

	client := resty.NewWithClient(&http.Client{
		Timeout: 10 * time.Second,
	})

	req := client.R()

	if webhook.Secret.Valid {
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		hmacSignature := generateHMACSHA256Secret(payloadBytes, webhook.Secret.String)
		req = req.SetHeader("x-hub-signature-256", fmt.Sprintf("sha256=%v", hmacSignature))
	}

	resp, err := req.SetBody(payload).Post(webhook.WebhookURL)
	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("webhook call failed with status code: %v and error: %v", resp.StatusCode(), resp.Error())
	}

	return nil
}

func generateHMACSHA256Secret(data []byte, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}
