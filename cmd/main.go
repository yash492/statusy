package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yash492/statusy/internal/applications"
	"github.com/yash492/statusy/internal/config"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	logger.Info("loading config")
	cfg := config.LoadConfig("")

	readDBConn := cfg.PostgresDB.ReadDB.String()
	writeDBConn := cfg.PostgresDB.WriteDB.String()

	logger.Info("connecting read database")
	readDB, err := pgxpool.New(ctx, readDBConn)
	if err != nil {
		logger.Error("failed to connect read database", slog.Any("err", err))
		os.Exit(1)
	}
	defer func() {
		logger.Info("closing read database connection")
		readDB.Close()
	}()

	logger.Info("connecting write database")
	writeDB, err := pgxpool.New(ctx, writeDBConn)
	if err != nil {
		logger.Error("failed to connect write database", slog.Any("err", err))
		os.Exit(1)
	}
	defer func() {
		logger.Info("closing write database connection")
		writeDB.Close()
	}()

	deps := applications.NewServerDeps(logger, readDB, writeDB)
	app := applications.NewServerApplication(deps)

	if err := app.Start(ctx, ":8081"); err != nil {
		logger.Error("server failed", slog.Any("err", err))
		os.Exit(1)
	}
}
