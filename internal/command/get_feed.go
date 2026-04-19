package command

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"time"

	"github.com/gorilla/feeds"
	"github.com/yash492/statusy/internal/domain/incidents"
	"github.com/yash492/statusy/internal/domain/scheduledmaintenance"
	"github.com/yash492/statusy/internal/domain/services"
)

type GetFeedCmd struct {
	lg                       *slog.Logger
	servicesRepo             services.Repository
	incidentsRepo            incidents.Repository
	scheduledMaintenanceRepo scheduledmaintenance.Repository
}

func NewGetFeedCmd(
	lg *slog.Logger,
	servicesRepo services.Repository,
	incidentsRepo incidents.Repository,
	scheduledMaintenanceRepo scheduledmaintenance.Repository,
) GetFeedCmd {
	return GetFeedCmd{
		lg:                       lg,
		servicesRepo:             servicesRepo,
		incidentsRepo:            incidentsRepo,
		scheduledMaintenanceRepo: scheduledMaintenanceRepo,
	}
}

type GetFeedParams struct {
	StatuspageSlug string
}

func (c GetFeedCmd) Execute(ctx context.Context, params GetFeedParams) (*feeds.Feed, error) {
	serviceResult, err := c.servicesRepo.GetBySlug(ctx, params.StatuspageSlug)
	if err != nil {
		c.lg.ErrorContext(ctx, "failed to get service by slug", slog.String("slug", params.StatuspageSlug), slog.Any("err", err))
		return nil, err
	}

	incResult, err := c.incidentsRepo.GetFeedByService(ctx, incidents.IncidentByServiceParams{
		ServiceID: serviceResult.ID,
		Limit:     50,
		Offset:    0,
	})
	if err != nil {
		c.lg.ErrorContext(ctx, "failed to get feed incidents", slog.String("slug", params.StatuspageSlug), slog.Any("err", err))
		return nil, err
	}

	smResult, err := c.scheduledMaintenanceRepo.GetFeedByService(ctx, scheduledmaintenance.ScheduledMaintenanceByServiceParams{
		ServiceID: serviceResult.ID,
		Limit:     50,
		Offset:    0,
	})
	if err != nil {
		c.lg.ErrorContext(ctx, "failed to get feed scheduled maintenances", slog.String("slug", params.StatuspageSlug), slog.Any("err", err))
		return nil, err
	}

	now := time.Now()
	feed := &feeds.Feed{
		Title:       fmt.Sprintf("%s Status", serviceResult.Name),
		Link:        &feeds.Link{Href: fmt.Sprintf("https://status.claude.com/")}, // TODO: update actual url if app has one
		Description: fmt.Sprintf("Status updates and incidents for %s", serviceResult.Name),
		Created:     now,
	}

	var items []*feeds.Item

	for _, inc := range incResult {
		description := fmt.Sprintf("Status: %s", inc.Status)
		if inc.AffectedComponents != "" {
			description += fmt.Sprintf("<br/>Affected Components: %s", inc.AffectedComponents)
		}

		items = append(items, &feeds.Item{
			Id:          fmt.Sprintf("incident-%d", inc.ID),
			Title:       fmt.Sprintf("Incident: %s", inc.Title),
			Link:        &feeds.Link{Href: inc.Link},
			Description: description,
			Created:     inc.ProviderCreatedAt,
		})
	}

	for _, sm := range smResult {
		description := fmt.Sprintf("Status: %s", sm.Status)
		if sm.AffectedComponents != "" {
			description += fmt.Sprintf("<br/>Affected Components: %s", sm.AffectedComponents)
		}

		items = append(items, &feeds.Item{
			Id:          fmt.Sprintf("maintenance-%d", sm.ID),
			Title:       fmt.Sprintf("Maintenance: %s", sm.Title),
			Link:        &feeds.Link{Href: sm.Link},
			Description: description,
			Created:     sm.ProviderCreatedAt,
		})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Created.After(items[j].Created)
	})

	if len(items) > 50 {
		items = items[:50]
	}

	feed.Items = items
	return feed, nil
}
