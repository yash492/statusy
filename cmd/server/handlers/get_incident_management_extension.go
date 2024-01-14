package handlers

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/pkg/api"
	"github.com/yash492/statusy/pkg/domain"
)

type incidentManagementExtensionResp struct {
	Squadcast squadcastHelperResp `json:"squadcast"`
	Pagerduty pagerdutyHelperResp `json:"pagerduty"`
}

type squadcastHelperResp struct {
	UUID         string `json:"uuid"`
	WebhookURL   string `json:"webhook_url"`
	IsConfigured bool   `json:"is_configured"`
}

type pagerdutyHelperResp struct {
	UUID         string `json:"uuid"`
	RoutingKey   string `json:"routing_key"`
	IsConfigured bool   `json:"is_configured"`
}

func GetIncidentManagementExtension(w http.ResponseWriter, r *http.Request) *api.Response {
	isSquadcastConfigured := true
	isPagerdutyConfigured := true

	squadcastExtension, err := domain.SquadcastExtension.Get()
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			isSquadcastConfigured = false

		} else {
			return api.Errorf(w, http.StatusInternalServerError, "cannot get squadcast extension configuration")
		}
	}

	pagerdutyExtension, err := domain.PagerdutyExtension.Get()
	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			isPagerdutyConfigured = false
		} else {
			return api.Errorf(w, http.StatusInternalServerError, "cannot get pagerduty extension configuration")
		}
	}
	squadcastUUID := ""
	if squadcastExtension.UUID != uuid.Nil {
		squadcastUUID = squadcastExtension.UUID.String()
	}

	pagerdutyUUID := ""
	if pagerdutyExtension.UUID != uuid.Nil {
		pagerdutyUUID = pagerdutyExtension.UUID.String()
	}

	return api.Send(w, http.StatusOK, incidentManagementExtensionResp{
		Squadcast: squadcastHelperResp{
			WebhookURL:   squadcastExtension.WebhookURL,
			IsConfigured: isSquadcastConfigured,
			UUID:         squadcastUUID,
		},
		Pagerduty: pagerdutyHelperResp{
			RoutingKey:   pagerdutyExtension.RoutingKey,
			IsConfigured: isPagerdutyConfigured,
			UUID:         pagerdutyUUID,
		},
	}, nil)
}
