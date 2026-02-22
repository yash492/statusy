package componentgroups

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yash492/statusy/internal/repository/componentgroups"
)

var _ componentgroups.ComponentsGroupRepository = &PostgresComponentGroupsRepository{}

type PostgresComponentGroupsRepository struct {
	lg      *slog.Logger
	writeDB *pgxpool.Pool
	readDB  *pgxpool.Pool
}

func NewPostgresComponentGroupsRepository(
	lg *slog.Logger,
	readDB *pgxpool.Pool,
	writeDB *pgxpool.Pool,
) *PostgresComponentGroupsRepository {
	return &PostgresComponentGroupsRepository{
		lg:      lg,
		writeDB: writeDB,
		readDB:  readDB,
	}
}
