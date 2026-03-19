package command

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/internal/domain/incidents"
	"github.com/yash492/statusy/internal/domain/services"
)

var ErrStatuspageNotFound = errors.New("statuspage not found")

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
	StatuspageSlug string
	PageNumber     int
	PageSize       int
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
}

func (c IncidentByStatuspageCmd) Execute(ctx context.Context, params IncidentByStatuspageParams) (IncidentByStatuspageResult, error) {
	slug := strings.TrimSpace(params.StatuspageSlug)
	if slug == "" {
		return IncidentByStatuspageResult{}, ErrStatuspageNotFound
	}

	service, err := c.ServicesRepo.GetBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.logger.WarnContext(ctx, "statuspage not found", slog.String("slug", slug))
			return IncidentByStatuspageResult{}, ErrStatuspageNotFound
		}

		c.logger.ErrorContext(ctx, "failed to fetch statuspage service", slog.String("slug", slug), slog.Any("err", err))
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
		ServiceID: service.ID,
		Limit:     pageSize,
		Offset:    offset,
	})
	if err != nil {
		c.logger.ErrorContext(ctx, "failed to fetch incidents by statuspage", slog.String("slug", slug), slog.Any("service_id", service.ID), slog.Any("err", err))
		return IncidentByStatuspageResult{}, err
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
		ServiceName: service.Title,
		ServiceSlug: service.Slug,
	}

	return result, nil
}
