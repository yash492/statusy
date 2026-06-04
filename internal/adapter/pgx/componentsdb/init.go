package componentsdb

import (
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	domaincomponents "github.com/yash492/statusy/internal/domain/components"
)

var _ domaincomponents.Repository = &PostgresComponentRepository{}

type PostgresComponentRepository struct {
	lg      *slog.Logger
	writeDB *pgxpool.Pool
	readDB  *pgxpool.Pool
}

func NewPostgresComponentRepository(
	lg *slog.Logger,
	readDB *pgxpool.Pool,
	writeDB *pgxpool.Pool,
) *PostgresComponentRepository {
	return &PostgresComponentRepository{
		lg:      lg,
		writeDB: writeDB,
		readDB:  readDB,
	}
}

type componentDto struct {
	ID               uint               `db:"id"`
	Name             string             `db:"name"`
	ProviderID       string             `db:"provider_id"`
	ServiceID        uint               `db:"service_id"`
	ComponentGroupID pgtype.Uint64      `db:"component_group_id"`
	CreatedAt        time.Time          `db:"created_at"`
	UpdatedAt        time.Time          `db:"updated_at"`
	DeletedAt        pgtype.Timestamptz `db:"deleted_at"`
}
