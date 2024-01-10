package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/yash492/statusy/pkg/api"
	"github.com/yash492/statusy/pkg/domain"
	"github.com/yash492/statusy/pkg/types"
)

type saveWebhookExtensionReq struct {
	WebhookURL string  `json:"webhook_url"`
	Secret     *string `json:"secret"`
	UUID       *string `json:"uuid"`
}

func (s *saveWebhookExtensionReq) Validate() error {
	if strings.TrimSpace(s.WebhookURL) == "" {
		return fmt.Errorf("%w: webhook url cannot be empty", api.ErrValidation)
	}
	if s.UUID != nil {
		if err := uuid.Validate(*s.UUID); err != nil {
			return fmt.Errorf("%w: invalid uuid format", api.ErrValidation)
		}
	}
	return nil
}

func SaveWebhookExtension(w http.ResponseWriter, r *http.Request) *api.Response {
	var req saveWebhookExtensionReq
	if err := api.Decode(r, &req); err != nil {
		if errors.Is(err, api.ErrValidation) {
			return api.Errorf(w, http.StatusBadRequest, err.Error())
		}
	}

	webhookUUID := uuid.New()
	if req.UUID != nil {
		webhookUUID = uuid.MustParse(*req.UUID)
	}

	if err := domain.WebhookExtension.Save(req.WebhookURL, req.Secret, webhookUUID); err != nil {
		return api.Errorf(w, http.StatusInternalServerError, "cannot save squadcast extension, err: %v", err.Error())
	}

	return api.Send(w, http.StatusOK, types.JSON{
		"msg": "squadcast extension is successfully saved",
	}, nil)
}
