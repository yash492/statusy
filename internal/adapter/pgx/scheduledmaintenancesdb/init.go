package scheduledmaintenancesdb

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresScheduledMaintenanceRepository struct {
	lg      *slog.Logger
	writeDB *pgxpool.Pool
	readDB  *pgxpool.Pool
}

func NewPostgresScheduledMaintenanceRepository(
	lg *slog.Logger,
	readDB *pgxpool.Pool,
	writeDB *pgxpool.Pool,
) *PostgresScheduledMaintenanceRepository {
	return &PostgresScheduledMaintenanceRepository{
		lg:      lg,
		writeDB: writeDB,
		readDB:  readDB,
	}
}
