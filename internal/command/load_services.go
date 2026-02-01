package command

import (
	"context"
	"log/slog"
	"os"

	"github.com/goccy/go-yaml"
	servicesRepo "github.com/yash492/statusy/internal/repository/services"
)

type LoadServicesCmd struct {
	lg          *slog.Logger
	serviceRepo servicesRepo.ServiceRepository
}

func (s *LoadServicesCmd) Execute(ctx context.Context) ([]servicesRepo.ServiceResult, error) {
	lg := s.lg
	servicesYamlBytes, err := os.ReadFile("./data/services.yaml")
	if err != nil {
		lg.Error("error fetching the services from yaml", slog.Any("err", err))
		return nil, err
	}

	var servicesYaml []servicesRepo.ServiceParams

	err = yaml.UnmarshalContext(ctx, servicesYamlBytes, &servicesYaml)
	if err != nil {
		return nil, err
	}

	err = s.serviceRepo.SaveAll(ctx, servicesYaml)
	if err != nil {
		return nil, err
	}

	services, err := s.serviceRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return services, nil

}
