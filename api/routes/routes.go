package routes

import (
	"context"
	"database/sql"
	"log"

	"cekkustomer.com/api/handlers/cekdata"
	"cekkustomer.com/api/handlers/files"
	"cekkustomer.com/configs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sethvargo/go-envconfig"
)

func NewRoutes(db *sql.DB) *gin.Engine {
	var config configs.AppConfiguration
	if err := envconfig.Process(context.Background(), &config); err != nil {
		log.Fatal(err.Error())
	}

	router := gin.Default()

	router.Use(cors.Default())
	router.Use(gin.Logger())
	router.HandleMethodNotAllowed = true

	v1 := router.Group("/v1")
	{
		check := v1.Group("/check")
		{
			check.GET("/match", cekdata.GetDPT(db))
			check.GET("/locate", cekdata.GetKec(db))
		}

		file := v1.Group("/files")
		{
			file.POST("/import", files.ImportExcel)
			file.GET("/read", files.ReadFile)
		}
	}

	return router
}
