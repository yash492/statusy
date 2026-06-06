package notificationsdb

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yash492/statusy/internal/domain/notifications"
)

type PostgresNotificationsRepository struct {
	lg      *slog.Logger
	readDB  *pgxpool.Pool
	writeDB *pgxpool.Pool
}

func NewPostgresNotificationsRepository(
	lg *slog.Logger,
	readDB *pgxpool.Pool,
	writeDB *pgxpool.Pool,
) *PostgresNotificationsRepository {
	return &PostgresNotificationsRepository{
		lg:      lg,
		readDB:  readDB,
		writeDB: writeDB,
	}
}

var _ notifications.NotificationsRepository = &PostgresNotificationsRepository{}
