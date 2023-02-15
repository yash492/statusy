package main

import (
	"backend/cmd"
	"backend/config"
	"backend/store"
	"log"
	"os"
)

func main() {

	err := config.Load()
	if err != nil {
		log.Fatalln(err)
	}

	db := store.InitDb()

	bytes, err := os.ReadFile("store/models.sql")

	if err != nil {
		log.Fatalln(err)
	}

	tx := db.Exec(string(bytes))

	if tx.Error != nil {
		log.Fatalln(err)
	}

	err = cmd.AddServicesAndComponentsToDb(db)
	if err != nil {
		log.Fatalln(err)
	}

}
