package main

import (
	"log"

	"cekkustomer.com/api/servers"
	"cekkustomer.com/pkg/database"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println(err.Error())
	} else {
		log.Println("connection successfully")
	}

	// connection for postgres
	db, err := database.Init()
	if err != nil {
		log.Println(err.Error())
		return
	} else {
		log.Println("connection pool successfully")
	}

	servers.Migrate(db)
	defer database.CloseDB()

	servers.Run(db)
}
