package handlers

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx"
	"github.com/yash492/statusy/pkg/api"
	"github.com/yash492/statusy/pkg/domain"
)

type chatOpsExtensionResp struct {
	Slack   chatOpsHelperResp `json:"slack"`
	MSTeams chatOpsHelperResp `json:"msteams"`
	Discord chatOpsHelperResp `json:"discord"`
}

type chatOpsHelperResp struct {
	IsConfigured bool   `json:"is_configured"`
	WebhookURL   string `json:"webhook_url"`
}

func GetChatopsExtension(w http.ResponseWriter, r *http.Request) *api.Response {
	chatopsExtension, err := domain.ChatopsExtension.Get()
	if errors.Is(err, pgx.ErrNoRows) {
		return api.Errorf(w, http.StatusInternalServerError, "cannot get squadcast extension configuration")
	}

	var resp chatOpsExtensionResp
	if len(chatopsExtension) < 1 {
		return api.Send(w, http.StatusOK, resp, nil)
	}

	for _, chatop := range chatopsExtension {
		chatOpDetail := chatOpsHelperResp{
			IsConfigured: true,
			WebhookURL:   chatop.WebhookURL,
		}

		if chatop.Type == "slack" {
			resp.Slack = chatOpDetail
		}

		if chatop.Type == "msteams" {
			resp.MSTeams = chatOpDetail
		}

		if chatop.Type == "discord" {
			resp.Discord = chatOpDetail
		}
	}

	return api.Send(w, http.StatusOK, resp, nil)

}
