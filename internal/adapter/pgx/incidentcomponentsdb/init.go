package incidentcomponentsdb

import (
"log/slog"

"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresIncidentComponentsRepository struct {
	lg      *slog.Logger
	writeDB *pgxpool.Pool
	readDB  *pgxpool.Pool
}

func NewPostgresIncidentComponentsRepository(
lg *slog.Logger,
readDB *pgxpool.Pool,
writeDB *pgxpool.Pool,
) *PostgresIncidentComponentsRepository {
	return &PostgresIncidentComponentsRepository{
		lg:      lg,
		writeDB: writeDB,
		readDB:  readDB,
	}
}
