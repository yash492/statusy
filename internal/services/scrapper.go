package applications

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Deps struct {
	Logger  *slog.Logger
	ReadDB  *pgxpool.Pool
	WriteDB *pgxpool.Pool
}

func StartScrapper(deps Deps) {

}
