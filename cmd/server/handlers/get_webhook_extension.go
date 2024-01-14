package handlers

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/pkg/api"
	"github.com/yash492/statusy/pkg/domain"
)

type getWebhookExtensionResp struct {
	UUID         string `json:"uuid"`
	IsConfigured bool   `json:"is_configured"`
	WebhookURL   string `json:"webhook_url"`
	Secret       string `json:"secret"`
}

func GetWebhookExtension(w http.ResponseWriter, r *http.Request) *api.Response {
	isConfigured := true

	webhookExtension, err := domain.WebhookExtension.Get()
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			isConfigured = false
		} else {
			return api.Errorf(w, http.StatusInternalServerError, "cannot get webhook extension configuration %v", err.Error())
		}
	}

	webhookUUID := ""
	if webhookExtension.UUID != uuid.Nil {
		webhookUUID = webhookExtension.UUID.String()
	}

	return api.Send(w, http.StatusOK, getWebhookExtensionResp{
		IsConfigured: isConfigured,
		WebhookURL:   webhookExtension.WebhookURL,
		Secret:       webhookExtension.Secret.String,
		UUID:         webhookUUID,
	}, nil)

}
