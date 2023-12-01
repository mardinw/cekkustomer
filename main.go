package main

import (
	"cekkustomer/api/servers"
	"cekkustomer/pkg/database"
	"log"

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

	defer database.CloseDB()

	servers.Run()
}
