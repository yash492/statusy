package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yash492/statusy/internal/applications"
	"github.com/yash492/statusy/internal/config"
	"github.com/yash492/statusy/schema"
	"golang.org/x/sync/errgroup"
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

	if err := schema.MigrateFs(writeDB, logger); err != nil {
		logger.Error("failed to run migrations", slog.Any("err", err))
		os.Exit(1)
	}

	deps := applications.NewServerDeps(logger, readDB, writeDB)
	serverApp := applications.NewServerApplication(deps)
	scrapperDeps := applications.NewScrapperDeps(logger, readDB, writeDB)
	scrapperApp := applications.NewScrapperApplication(scrapperDeps)

	errGroup := new(errgroup.Group)

	errGroup.Go(func() error {
		return serverApp.Start(ctx, fmt.Sprintf(":%d", cfg.ServerPort))
	})

	errGroup.Go(func() error {
		return scrapperApp.Start(ctx, cfg.ScrappingInterval)
	})

	if err := errGroup.Wait(); err != nil {
		logger.Error("server or scrapper stopped with error", slog.Any("err", err))
		os.Exit(1)
	}

}
