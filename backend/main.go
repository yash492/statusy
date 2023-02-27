package main

import (
	"backend/cmd"
	"backend/config"
	"backend/store"
	"log"
	"os"
)

func main() {

	err := config.Load("./config.yaml")
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

	err = cmd.Init()
	if err != nil {
		log.Fatalln(err)
	}

}
