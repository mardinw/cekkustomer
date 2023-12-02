package routes

import (
	"context"
	"database/sql"
	"log"

	"cekkustomer.com/api/handlers/cekdata"
	"cekkustomer.com/configs"

	"github.com/gin-gonic/gin"
	"github.com/sethvargo/go-envconfig"
)

func NewRoutes(db *sql.DB) *gin.Engine {
	var config configs.AppConfiguration
	if err := envconfig.Process(context.Background(), &config); err != nil {
		log.Fatal(err.Error())
	}

	router := gin.Default()

	router.Use(gin.Logger())
	router.HandleMethodNotAllowed = true

	v1 := router.Group("/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})

		check := v1.Group("/check")
		{
			check.GET("/match", cekdata.GetDPT(db))
		}
	}

	return router
}
