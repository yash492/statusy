package command

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/yash492/statusy/internal/domain/components"
	"github.com/yash492/statusy/internal/domain/services"
)

type GetServiceComponentsCmd struct {
	logger              *slog.Logger
	servicesRepo        services.Repository
	componentsRepo      components.Repository
	componentGroupsRepo components.GroupRepository
}

func NewGetServiceComponentsCmd(
	logger *slog.Logger,
	servicesRepo services.Repository,
	componentsRepo components.Repository,
	componentGroupsRepo components.GroupRepository,
) GetServiceComponentsCmd {
	return GetServiceComponentsCmd{
		logger:              logger,
		servicesRepo:        servicesRepo,
		componentsRepo:      componentsRepo,
		componentGroupsRepo: componentGroupsRepo,
	}
}

type GetServiceComponentsParams struct {
	ServiceSlug string
}

type GetServiceComponentsResult struct {
	ServiceID           uint
	ServiceName         string
	ServiceSlug         string
	GroupedComponents   []GroupedComponentResult
	UngroupedComponents []ComponentResult
}

type GroupedComponentResult struct {
	ID         uint
	Name       string
	ProviderID string
	Components []ComponentResult
}

type ComponentResult struct {
	ID         uint
	Name       string
	ProviderID string
}

func (c GetServiceComponentsCmd) Execute(ctx context.Context, params GetServiceComponentsParams) (GetServiceComponentsResult, error) {
	slug := strings.TrimSpace(params.ServiceSlug)
	if slug == "" {
		return GetServiceComponentsResult{}, ErrStatuspageNotFound
	}

	service, err := c.servicesRepo.GetBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.logger.WarnContext(ctx, "service not found by slug", slog.String("slug", slug))
			return GetServiceComponentsResult{}, ErrStatuspageNotFound
		}
		c.logger.ErrorContext(ctx, "failed to fetch service", slog.String("slug", slug), slog.Any("err", err))
		return GetServiceComponentsResult{}, err
	}

	groups, err := c.componentGroupsRepo.GetByServiceID(ctx, service.ID)
	if err != nil {
		c.logger.ErrorContext(ctx, "failed to fetch component groups", slog.Uint64("service_id", uint64(service.ID)), slog.Any("err", err))
		return GetServiceComponentsResult{}, err
	}

	comps, err := c.componentsRepo.GetByServiceID(ctx, service.ID)
	if err != nil {
		c.logger.ErrorContext(ctx, "failed to fetch components", slog.Uint64("service_id", uint64(service.ID)), slog.Any("err", err))
		return GetServiceComponentsResult{}, err
	}

	// Index components by component group ID
	groupedCompsMap := make(map[uint][]ComponentResult)
	var ungroupedComps []ComponentResult

	for _, comp := range comps {
		cgID, hasGroup := comp.ComponentGroupID.Get()
		compRes := ComponentResult{
			ID:         comp.ID,
			Name:       comp.Name,
			ProviderID: comp.ProviderID,
		}
		if hasGroup {
			groupedCompsMap[cgID] = append(groupedCompsMap[cgID], compRes)
		} else {
			ungroupedComps = append(ungroupedComps, compRes)
		}
	}

	groupedResults := make([]GroupedComponentResult, 0, len(groups))
	for _, group := range groups {
		gComps := groupedCompsMap[group.ID]
		if gComps == nil {
			gComps = []ComponentResult{}
		}
		groupedResults = append(groupedResults, GroupedComponentResult{
			ID:         group.ID,
			Name:       group.Name,
			ProviderID: group.ProviderID,
			Components: gComps,
		})
	}

	// Double check: if ungroupedComps is nil, make it an empty slice
	if ungroupedComps == nil {
		ungroupedComps = []ComponentResult{}
	}

	return GetServiceComponentsResult{
		ServiceID:           service.ID,
		ServiceName:         service.Name,
		ServiceSlug:         service.Slug,
		GroupedComponents:   groupedResults,
		UngroupedComponents: ungroupedComps,
	}, nil
}
