package store

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type dbCreds struct {
	host     string
	password string
	dbname   string
	user     string
	sslmode  string
	port     int
}

func InitDb() *gorm.DB {
	dsn := dbCreds{
		host:     "localhost",
		port:     5432,
		password: "root",
		dbname:   "statusy",
		user:     "yash",
		sslmode:  "disable",
	}

	db, err := gorm.Open(postgres.Open(dsn.string()), &gorm.Config{})

	if err != nil {
		log.Fatalln("could not initialise db ", err)
		return nil
	}

	return db
}
func (d dbCreds) string() string {
	return fmt.Sprintf("host=%v password=%v dbname=%v user=%v sslmode=%v port=%v", d.host, d.password, d.dbname, d.user, d.sslmode, d.port)
}

func InitDbEnv(db *gorm.DB) *Db {
	return &Db{db}
}
