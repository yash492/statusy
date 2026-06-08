package viewsdb

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yash492/statusy/internal/domain/views"
)

var _ views.Repository = &PostgresViewsRepository{}

type PostgresViewsRepository struct {
	lg      *slog.Logger
	readDB  *pgxpool.Pool
	writeDB *pgxpool.Pool
}

func NewPostgresViewsRepository(
	lg *slog.Logger,
	readDB *pgxpool.Pool,
	writeDB *pgxpool.Pool,
) *PostgresViewsRepository {
	return &PostgresViewsRepository{
		lg:      lg,
		readDB:  readDB,
		writeDB: writeDB,
	}
}
