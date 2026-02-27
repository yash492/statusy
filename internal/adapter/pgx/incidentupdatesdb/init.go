package incidentupdatesdb

import (
"log/slog"

"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresIncidentUpdatesRepository struct {
	lg      *slog.Logger
	writeDB *pgxpool.Pool
	readDB  *pgxpool.Pool
}

func NewPostgresIncidentUpdatesRepository(
lg *slog.Logger,
readDB *pgxpool.Pool,
writeDB *pgxpool.Pool,
) *PostgresIncidentUpdatesRepository {
	return &PostgresIncidentUpdatesRepository{
		lg:      lg,
		writeDB: writeDB,
		readDB:  readDB,
	}
}
