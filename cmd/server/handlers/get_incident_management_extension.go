package handlers

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx"
	"github.com/yash492/statusy/pkg/api"
	"github.com/yash492/statusy/pkg/domain"
)

type incidentManagementExtensionResp struct {
	Squadcast squadcastHelperResp `json:"squadcast"`
	Pageduty  pagerdutyHelperResp `json:"pageduty"`
}

type squadcastHelperResp struct {
	WebhookURL   string `json:"webhook_url"`
	IsConfigured bool   `json:"is_configured"`
}

type pagerdutyHelperResp struct {
	RoutingKey   string `json:"routing_key"`
	IsConfigured bool   `json:"is_configured"`
}

func GetIncidentManagementExtension(w http.ResponseWriter, r *http.Request) *api.Response {
	isSquadcastConfigured := true
	isPagerdutyConfigured := true

	squadcastExtension, err := domain.SquadcastExtension.Get()
	if errors.Is(err, pgx.ErrNoRows) {
		isSquadcastConfigured = true
	} else {
		return api.Errorf(w, http.StatusInternalServerError, "cannot get squadcast extension configuration")
	}

	pagerdutyExtension, err := domain.PagerdutyExtension.Get()
	if errors.Is(err, pgx.ErrNoRows) {
		isPagerdutyConfigured = true
	} else {
		return api.Errorf(w, http.StatusInternalServerError, "cannot get pagerduty extension configuration")
	}

	return api.Send(w, http.StatusOK, incidentManagementExtensionResp{
		Squadcast: squadcastHelperResp{
			WebhookURL:   squadcastExtension.WebhookURL,
			IsConfigured: isSquadcastConfigured,
		},
		Pageduty: pagerdutyHelperResp{
			RoutingKey:   pagerdutyExtension.RoutingKey,
			IsConfigured: isPagerdutyConfigured,
		},
	}, nil)
}
