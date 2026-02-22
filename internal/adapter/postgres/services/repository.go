package services

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yash492/statusy/internal/repository/services"
)

var _ services.ServiceRepository = &PostgresServiceRepository{}

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
