package incidentsdb

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresIncidentRepository struct {
	lg      *slog.Logger
	writeDB *pgxpool.Pool
	readDB  *pgxpool.Pool
}

func NewPostgresIncidentRepository(
	lg *slog.Logger,
	readDB *pgxpool.Pool,
	writeDB *pgxpool.Pool,
) *PostgresIncidentRepository {
	return &PostgresIncidentRepository{
		lg:      lg,
		writeDB: writeDB,
		readDB:  readDB,
	}
}
