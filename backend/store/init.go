package store

import (
	"backend/config"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type dbCreds struct {
	host     string
	password string
	dbname   string
	user     string
	sslmode  string
	port     int
}

var dbEnv *gorm.DB

func InitDb() *gorm.DB {
	dsn := dbCreds{
		host:     config.Env.Db.Host,
		port:     config.Env.Db.Port,
		password: config.Env.Db.Password,
		dbname:   config.Env.Db.Dbname,
		user:     config.Env.Db.User,
		sslmode:  config.Env.Db.SslMode,
	}

	db, err := gorm.Open(postgres.Open(dsn.string()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalln("could not initialise db ", err)
		return nil
	}
	dbEnv = db
	return db
}
func (d dbCreds) string() string {
	return fmt.Sprintf("host=%v password=%v dbname=%v user=%v sslmode=%v port=%v", d.host, d.password, d.dbname, d.user, d.sslmode, d.port)
}

func InitDbEnv() *Db {
	return &Db{dbEnv}
}
