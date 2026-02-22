package components

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yash492/statusy/internal/repository/components"
)

var _ components.ComponentRepository = &PostgresComponentRepository{}

type PostgresComponentRepository struct {
	lg      *slog.Logger
	writeDB *pgxpool.Pool
	readDB  *pgxpool.Pool
}

func NewPostgresComponentRepository(
	lg *slog.Logger,
	readDB *pgxpool.Pool,
	writeDB *pgxpool.Pool,
) *PostgresComponentRepository {
	return &PostgresComponentRepository{
		lg:      lg,
		writeDB: writeDB,
		readDB:  readDB,
	}
}
