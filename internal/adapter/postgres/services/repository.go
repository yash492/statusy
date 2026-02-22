package services

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	domainservices "github.com/yash492/statusy/internal/domain/services"
)

var _ domainservices.Repository = &PostgresServiceRepository{}

type PostgresServiceRepository struct {
	lg      *slog.Logger
	writeDB *pgxpool.Pool
	readDB  *pgxpool.Pool
}

func NewPostgresServiceRepository(
	lg *slog.Logger,
	readDB *pgxpool.Pool,
	writeDB *pgxpool.Pool,
) *PostgresServiceRepository {
	return &PostgresServiceRepository{
		lg:      lg,
		writeDB: writeDB,
		readDB:  readDB,
	}
}
