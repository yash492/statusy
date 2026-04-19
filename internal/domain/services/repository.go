package services

import (
	"context"
)

type ServiceResult struct {
	ID   uint
	Name string
	Slug string
}

type ServiceParams struct {
	Name string
	Slug string
}

type Repository interface {
	SaveAll(ctx context.Context, services []ServiceParams) ([]ServiceResult, error)
	SearchByName(ctx context.Context, name string) ([]ServiceResult, error)
	GetBySlug(ctx context.Context, slug string) (ServiceResult, error)
}
