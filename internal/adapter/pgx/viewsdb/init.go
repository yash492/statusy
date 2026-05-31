package viewsdb

import (
	"log/slog"
	"time"

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

// viewDto maps a row from the views table.
type viewDto struct {
	ID          uint       `db:"id"`
	Name        string     `db:"name"`
	PublicID    string     `db:"public_id"`
	Description string     `db:"description"`
	IsDefault   bool       `db:"is_default"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
}

// viewServiceDto maps a row joining view_services → services.
type viewServiceDto struct {
	ID                   uint       `db:"id"`
	Name                 string     `db:"name"`
	Slug                 string     `db:"slug"`
	IncludeAllComponents bool       `db:"include_all_components"`
	UpdatedAt            time.Time  `db:"updated_at"`
	DeletedAt            *time.Time `db:"deleted_at"`
}

// viewServiceFullDto maps a full row from the view_services table.
type viewServiceFullDto struct {
	ID                   uint       `db:"id"`
	ViewID               uint       `db:"view_id"`
	ServiceID            uint       `db:"service_id"`
	IncludeAllComponents bool       `db:"include_all_components"`
	CreatedAt            time.Time  `db:"created_at"`
	UpdatedAt            time.Time  `db:"updated_at"`
	DeletedAt            *time.Time `db:"deleted_at"`
}
