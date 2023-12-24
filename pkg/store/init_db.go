package store

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yash492/statusy/pkg/config"
)

type db struct {
	pgConn *pgxpool.Pool
}

var dbConn db

func New() error {

	config, err := pgxpool.ParseConfig(dbConnectionStr())
	if err != nil {
		log.Fatalln(err)
	}

	conn, err := pgxpool.NewWithConfig(context.Background(), config)

	if err != nil {
		log.Fatalln(err)
	}

	if err = conn.Ping(context.Background()); err != nil {
		log.Fatalln(err)
	}

	dbConn = db{
		pgConn: conn,
	}

	return nil
}

func InitDbVar() db {
	return dbConn
}

func dbConnectionStr() string {
	str := fmt.Sprintf("postgres://%s:%s@%s:5432/%s", config.PGUser, config.PGPassword, config.PGHost, config.PGDatabaseName)
	return str
}
