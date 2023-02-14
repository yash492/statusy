package main

import (
	"backend/cmd"
	"backend/store"
	"log"
	"os"
)

func main() {
	db := store.InitDb()

	bytes, err := os.ReadFile("store/models.sql")

	if err != nil {
		log.Fatalln(err)
	}

	tx := db.Exec(string(bytes))

	if tx.Error != nil {
		log.Fatalln(err)
	}

	err = cmd.AddServicesToDb(db)
	if err != nil {
		log.Fatalln(err)
	}

	err = cmd.AddComponentsToDb(db)
	if err != nil {
		log.Fatalln(err)
	}

}
