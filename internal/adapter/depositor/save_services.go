package depositor

import (
	"context"
	"log"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yash492/statusy/internal/config"
	"github.com/yash492/statusy/sqlc/db"
)

func SaveServices() map[string]any {

	cfg := config.LoadConfig("")
	yamlFile, err := os.ReadFile("data/services.yaml")
	if err != nil {
		log.Fatalln(err)
	}

	var services []map[string]any

	err = yaml.Unmarshal(yamlFile, &services)
	if err != nil {
		log.Fatalln(err)
	}
	writeDBPool, err := pgxpool.New(context.Background(), cfg.PostgresDB.WriteDB.String())

	if err != nil {
		log.Fatal("unable to establish the connection with write db", err)
		os.Exit(1)
	}

	queries := db.New(writeDBPool)

	var sericesParams []db.CreateServicesParams

	for _, service := range services {
		sericesParams = append(sericesParams, db.CreateServicesParams{
			Name: service["name"].(string),
			Slug: service["name"].(string),
			ProviderType: service["name"].(string),
			IncidentsUrl: service["name"].(string),
			
			
		})
	}

	queries.CreateServices(context.Background())

	return services

}
