package command

import (
	"context"
	"log/slog"

	"github.com/goccy/go-yaml"
	domainservices "github.com/yash492/statusy/internal/domain/services"
)

type LoadServicesCmd struct {
	lg                *slog.Logger
	servicesYamlBytes []byte
	serviceRepo       domainservices.Repository
	
}

func NewLoadServiceCmd(
	lg *slog.Logger,
	serviceYamlBytes []byte,
	serviceRepo domainservices.Repository) *LoadServicesCmd {
	return &LoadServicesCmd{
		lg:                lg.With("cmd", "LoadServicesCmd"),
		servicesYamlBytes: serviceYamlBytes,
		serviceRepo:       serviceRepo,
	}
}

func (s *LoadServicesCmd) Execute(ctx context.Context) ([]domainservices.ServiceResult, error) {
	lg := s.lg

	var servicesYaml []domainservices.ServiceParams

	err := yaml.UnmarshalContext(ctx, s.servicesYamlBytes, &servicesYaml)
	if err != nil {
		lg.ErrorContext(ctx, "unable to marshal services yaml", slog.Any("err", err))
		return nil, err
	}

	services, err := s.serviceRepo.SaveAll(ctx, servicesYaml)
	if err != nil {
		lg.ErrorContext(ctx, "unable to save all services", slog.Any("err", err))
		return nil, err
	}

	return services, nil

}
