package main

import (
	"cekkustomer/api/servers"
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

	servers.Run()
}
