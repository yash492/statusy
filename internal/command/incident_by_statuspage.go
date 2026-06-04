package command

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/yash492/statusy/internal/common/apperrors"
	"github.com/yash492/statusy/internal/domain/incidents"
	"github.com/yash492/statusy/internal/domain/services"
)

type IncidentByStatuspageCmd struct {
	logger        *slog.Logger
	ServicesRepo  services.Repository
	incidentsRepo incidents.Repository
}

func NewIncidentByStatuspageCmd(
	logger *slog.Logger,
	servicesRepo services.Repository,
	incidentsRepo incidents.Repository,
) IncidentByStatuspageCmd {
	return IncidentByStatuspageCmd{
		logger:        logger,
		ServicesRepo:  servicesRepo,
		incidentsRepo: incidentsRepo,
	}
}

type IncidentByStatuspageParams struct {
	StatuspageSlug    string
	ComponentIDs      []int
	ComponentGroupIDs []int
	PageNumber        int
	PageSize          int
}

type IncidentByStatuspageIncident struct {
	ID                uint
	Title             string
	Status            string
	ProviderCreatedAt time.Time
	Link              string
}

type IncidentByStatuspageResult struct {
	Incidents   []IncidentByStatuspageIncident
	ServiceName string
	ServiceSlug string
	ServiceUrl  string
	TotalCount  int
	ServiceID   uint
}

func (c IncidentByStatuspageCmd) Execute(ctx context.Context, params IncidentByStatuspageParams) (IncidentByStatuspageResult, error) {
	slug := strings.TrimSpace(params.StatuspageSlug)
	if slug == "" {
		return IncidentByStatuspageResult{}, apperrors.InvalidInputError("slug cannot be empty", fmt.Errorf("slug cannot be empty"))
	}

	service, err := c.ServicesRepo.GetBySlug(ctx, slug)
	if err != nil {
		return IncidentByStatuspageResult{}, err
	}

	pageNumber := params.PageNumber
	if pageNumber <= 0 {
		pageNumber = 1
	}

	pageSize := params.PageSize

	//Keeping it hardcoded for now
	if pageSize != 10 {
		pageSize = 10
	}

	offset := (pageNumber - 1) * pageSize

	incidentRows, err := c.incidentsRepo.GetByService(ctx, incidents.IncidentByServiceParams{
		ServiceID:         service.ID,
		ComponentIDs:      params.ComponentIDs,
		ComponentGroupIDs: params.ComponentGroupIDs,
		Limit:             pageSize,
		Offset:            offset,
	})
	if err != nil {
		return IncidentByStatuspageResult{}, err
	}

	totalCount := 0
	if len(incidentRows) > 0 {
		totalCount = int(incidentRows[0].TotalCount)
	}

	incidents := make([]IncidentByStatuspageIncident, 0, len(incidentRows))
	for _, incident := range incidentRows {
		incidents = append(incidents, IncidentByStatuspageIncident{
			ID:                incident.ID,
			Title:             incident.Title,
			Status:            incident.Status,
			ProviderCreatedAt: incident.ProviderCreatedAt,
			Link:              incident.Link,
		})
	}

	result := IncidentByStatuspageResult{
		Incidents:   incidents,
		ServiceName: service.Name,
		ServiceSlug: service.Slug,
		ServiceUrl:  service.URL,
		TotalCount:  totalCount,
		ServiceID:   service.ID,
	}

	return result, nil
}
