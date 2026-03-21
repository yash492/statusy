package scheduledmaintenanceupdatesdb

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresScheduledMaintenanceUpdatesRepository struct {
	lg      *slog.Logger
	writeDB *pgxpool.Pool
	readDB  *pgxpool.Pool
}

func NewPostgresScheduledMaintenanceUpdatesRepository(
	lg *slog.Logger,
	readDB *pgxpool.Pool,
	writeDB *pgxpool.Pool,
) *PostgresScheduledMaintenanceUpdatesRepository {
	return &PostgresScheduledMaintenanceUpdatesRepository{
		lg:      lg,
		writeDB: writeDB,
		readDB:  readDB,
	}
}
