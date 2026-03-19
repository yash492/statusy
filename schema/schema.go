package schema

import (
	"embed"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var embedFS embed.FS

func MigrateFs(dbPool *pgxpool.Pool, logger *slog.Logger) error {
	gooseLogger := gooseLogger{
		logger: logger,
	}
	goose.SetLogger(&gooseLogger)
	goose.SetBaseFS(embedFS)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	db := stdlib.OpenDBFromPool(dbPool)
	if err := goose.Up(db, "."); err != nil {
		return err
	}
	return nil
}

type gooseLogger struct {
	logger *slog.Logger
}

func (g *gooseLogger) Fatalf(format string, v ...any) {
	g.logger.Error(fmt.Sprintf(format, v...))

}
func (g *gooseLogger) Printf(format string, v ...any) {
	g.logger.Info(fmt.Sprintf(format, v...))
}
