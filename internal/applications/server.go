package applications

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ServerDeps struct {
	lg      *slog.Logger
	readDB  *pgxpool.Pool
	writeDB *pgxpool.Pool
}

func NewServerDeps(
	lg *slog.Logger,
	readDB *pgxpool.Pool,
	writeDB *pgxpool.Pool,
) ServerDeps {
	return ServerDeps{
		lg:      lg,
		readDB:  readDB,
		writeDB: writeDB,
	}
}

func (s ScrapperDeps) Start() {
	
}
