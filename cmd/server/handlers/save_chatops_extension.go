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

type saveChatOpsExtensionReq struct {
	Type       string  `json:"type"`
	WebhookURL string  `json:"webhook_url"`
	UUID       *string `json:"uuid"`
}

var validChatOpsTypes = map[string]bool{
	"slack":   true,
	"msteams": true,
	"discord": true,
}

func (c *saveChatOpsExtensionReq) Validate() error {
	if !validChatOpsTypes[c.Type] {
		return fmt.Errorf("%w: invalid chatops type", api.ErrValidation)
	}

	if strings.TrimSpace(c.WebhookURL) == "" {
		return fmt.Errorf("%w: webhook url cannot be empty", api.ErrValidation)
	}

	if c.UUID != nil {
		if err := uuid.Validate(*c.UUID); err != nil {
			return fmt.Errorf("%w: invalid uuid format", api.ErrValidation)
		}
	}

	return nil
}

func SaveChatOpsExtension(w http.ResponseWriter, r *http.Request) *api.Response {
	var req saveChatOpsExtensionReq
	if err := api.Decode(r, &req); err != nil {
		if errors.Is(err, api.ErrValidation) {
			return api.Errorf(w, http.StatusBadRequest, err.Error())
		}
	}

	chatOpsUUID := uuid.New()
	if req.UUID != nil {
		chatOpsUUID = uuid.MustParse(*req.UUID)
	}

	if err := domain.ChatopsExtension.Save(req.Type, req.WebhookURL, chatOpsUUID); err != nil {
		return api.Errorf(w, http.StatusInternalServerError, "cannot save %v extension", req.Type)
	}

	return api.Send(w, http.StatusOK, types.JSON{
		"msg": fmt.Sprintf("%v extension is successfully saved", req.Type),
	}, nil)
}
