package command

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	TestDb *pgxpool.Pool
}

func (t *TestSuite) SetupSuite() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", "statusy", "password", "localhost", 5432, "statusy")

	writeDBPool, err := pgxpool.New(context.Background(), connString)

	if err != nil {
		logger.Error("unable to establish the connection with write db", err)
		os.Exit(1)
	}

	t.TestDb = writeDBPool
}
