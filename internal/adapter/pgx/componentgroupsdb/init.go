package componentgroupsdb

import (
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
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

type componentGroupDto struct {
	ID         uint             `db:"id"`
	Name       string           `db:"name"`
	ProviderID string           `db:"provider_id"`
	ServiceID  uint             `db:"service_id"`
	CreatedAt  time.Time        `db:"created_at"`
	UpdatedAt  time.Time        `db:"updated_at"`
	DeletedAt  pgtype.Timestamp `db:"deleted_at"`
}
