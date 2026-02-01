package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yash492/statusy/internal/config"
	"github.com/yash492/statusy/internal/services"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	cfg := config.LoadConfig("")

	writeDBPool, err := pgxpool.New(context.Background(), cfg.PostgresDB.WriteDB.String())
	if err != nil {
		logger.Error("unable to establish the connection with write db", err)
		os.Exit(1)
	}

	readDBPool, err := pgxpool.New(context.Background(), cfg.PostgresDB.ReadDB.String())
	if err != nil {
		logger.Error("unable to establish the connection with write db", err)
		os.Exit(1)
	}

	deps := services.Deps{
		Logger:  logger,
		ReadDB:  readDBPool,
		WriteDB: writeDBPool,
	}
}
