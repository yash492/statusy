package config

import "github.com/knadh/koanf/v2"

var (
	PGHost              string
	PGPassword          string
	PGUser              string
	PGPort              int
	PGDatabaseName      string
	ScrapIntervalInMins int
)

func initVariables(k *koanf.Koanf) {

	PGHost = k.String("db.host")
	PGPassword = k.String("db.password")
	PGUser = k.String("db.user")
	PGPort = k.Int("db.port")
	PGDatabaseName = k.String("db.database_name")

	ScrapIntervalInMins = k.Int("app.scrap_interval")
	if ScrapIntervalInMins == 0 {
		ScrapIntervalInMins = 5
	}
}
