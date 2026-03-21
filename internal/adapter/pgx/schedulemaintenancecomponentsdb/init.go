package schedulemaintenancecomponentsdb

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresScheduleMaintenanceComponentsRepository struct {
	lg      *slog.Logger
	writeDB *pgxpool.Pool
	readDB  *pgxpool.Pool
}

func NewPostgresScheduleMaintenanceComponentsRepository(
	lg *slog.Logger,
	readDB *pgxpool.Pool,
	writeDB *pgxpool.Pool,
) *PostgresScheduleMaintenanceComponentsRepository {
	return &PostgresScheduleMaintenanceComponentsRepository{
		lg:      lg,
		writeDB: writeDB,
		readDB:  readDB,
	}
}
