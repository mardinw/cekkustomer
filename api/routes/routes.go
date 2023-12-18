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
			check.GET("/locate", cekdata.GetKec(db))
			check.GET("/look/:foldername/:filename", cekdata.CheckDPT(db))
		}

		file := v1.Group("/files")
		{
			file.POST("/import", files.ImportExcel(db))
			file.GET("/export/:foldername/:filename", files.ExportMatchExcel(db))
			file.GET("/read/:foldername/:filename", files.ReadFile)
			file.GET("/list/:folder", files.GetListFolder)
			file.GET("/download/:filename", files.DownloadSampleXlsx)
			file.DELETE("/:foldername/:filename", files.DeleteFile(db))
		}
	}

	return router
}
