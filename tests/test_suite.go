package tests

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
	"github.com/yash492/statusy/internal/config"
)

type TestSuite struct {
	suite.Suite
	TestDb *pgxpool.Pool
	Cfg    config.Config
}

func (t *TestSuite) SetupSuite() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	cfg := config.LoadConfig("")

	writeDBPool, err := pgxpool.New(context.Background(), cfg.PostgresDB.WriteDB.String())

	if err != nil {
		logger.Error("unable to establish the connection with write db", err)
		os.Exit(1)
	}

	t.Cfg = cfg
	t.TestDb = writeDBPool
}
