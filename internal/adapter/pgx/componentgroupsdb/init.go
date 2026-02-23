package componentgroupsdb

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	domaincomponents "github.com/yash492/statusy/internal/domain/components"
)

var _ domaincomponents.GroupRepository = &PostgresComponentGroupsRepository{}

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
