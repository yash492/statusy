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

type savePagerdutyExtensionReq struct {
	RoutingKey string  `json:"routing_key"`
	UUID       *string `json:"uuid"`
}

func (s *savePagerdutyExtensionReq) Validate() error {
	if strings.TrimSpace(s.RoutingKey) == "" {
		return fmt.Errorf("%w: routing key cannot be empty", api.ErrValidation)
	}

	if s.UUID != nil {
		if err := uuid.Validate(*s.UUID); err != nil {
			return fmt.Errorf("%w: invalid uuid format", api.ErrValidation)
		}
	}
	return nil
}

func SavePagerdutyExtension(w http.ResponseWriter, r *http.Request) *api.Response {
	var req savePagerdutyExtensionReq
	if err := api.Decode(r, &req); err != nil {
		if errors.Is(err, api.ErrValidation) {
			return api.Errorf(w, http.StatusBadRequest, err.Error())
		}
	}

	pagerdutyUUID := uuid.New()
	if req.UUID != nil {
		pagerdutyUUID = uuid.MustParse(*req.UUID)
	}

	if err := domain.PagerdutyExtension.Save(req.RoutingKey, pagerdutyUUID); err != nil {
		return api.Errorf(w, http.StatusInternalServerError, "cannot save pagerduty extension, err: %v", err.Error())
	}

	return api.Send(w, http.StatusOK, types.JSON{
		"msg": "pagerduty extension is successfully saved",
	}, nil)
}
