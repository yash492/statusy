package services

import (
	"context"

	domainservices "github.com/yash492/statusy/internal/domain/services"
)

func (s *PostgresServiceRepository) GetAll(ctx context.Context) ([]domainservices.ServiceResult, error) {
	return nil, nil
}
