package command

import (
	"context"
	"log/slog"

	"github.com/goccy/go-yaml"
	servicesRepo "github.com/yash492/statusy/internal/repository/services"
)

type LoadServicesCmd struct {
	lg                *slog.Logger
	servicesYamlBytes []byte
	serviceRepo       servicesRepo.ServiceRepository
	
}

func NewLoadServiceCmd(
	lg *slog.Logger,
	serviceYamlBytes []byte,
	serviceRepo servicesRepo.ServiceRepository) *LoadServicesCmd {
	return &LoadServicesCmd{
		lg:                lg.With("cmd", "LoadServicesCmd"),
		servicesYamlBytes: serviceYamlBytes,
		serviceRepo:       serviceRepo,
	}
}

func (s *LoadServicesCmd) Execute(ctx context.Context) ([]servicesRepo.ServiceResult, error) {
	lg := s.lg

	var servicesYaml []servicesRepo.ServiceParams

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
