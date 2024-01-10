package handlers

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx"
	"github.com/yash492/statusy/pkg/api"
	"github.com/yash492/statusy/pkg/domain"
)

type getWebhookExtensionResp struct {
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
			return api.Errorf(w, http.StatusInternalServerError, "cannot get webhook extension configuration")
		}
	}

	return api.Send(w, http.StatusOK, getWebhookExtensionResp{
		IsConfigured: isConfigured,
		WebhookURL:   webhookExtension.WebhookURL,
		Secret:       webhookExtension.Secret.String,
	}, nil)

}
