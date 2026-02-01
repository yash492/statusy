package postgres

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yash492/statusy/internal/repository/services"
	"github.com/yash492/statusy/sqlc/db"
)

var _ services.ServiceRepository = &PostgresServiceRepository{}

type PostgresServiceRepository struct {
	lg      *slog.Logger
	writeDB *pgxpool.Pool
	readDB  *pgxpool.Pool
}

func NewPostgresServiceRepository(
	lg *slog.Logger,
	readDB *pgxpool.Pool,
	writeDB *pgxpool.Pool,
) *PostgresServiceRepository {
	return &PostgresServiceRepository{
		lg:      lg,
		writeDB: writeDB,
		readDB:  readDB,
	}
}

func (p *PostgresServiceRepository) write(fn func(q *db.Queries) error) error {
	queries := db.New(p.writeDB)
	err := fn(queries)
	return err
}

func (p *PostgresServiceRepository) read(fn func(q *db.Queries) error) error {
	queries := db.New(p.writeDB)
	err := fn(queries)
	return err
}
