package routes

import (
	"context"
	"database/sql"
	"log"
	"time"

	"cekkustomer.com/api/handlers/auth"
	"cekkustomer.com/api/handlers/cekdata"
	"cekkustomer.com/api/handlers/files"
	"cekkustomer.com/api/middlewares"
	"cekkustomer.com/configs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sethvargo/go-envconfig"
)

func NewRoutes(db *sql.DB) *gin.Engine {
	var config configs.AppConfiguration
	var ttlDynmo configs.AwsDynTblConfig

	if err := envconfig.Process(context.Background(), &config); err != nil {
		log.Fatal(err.Error())
	}

	if err := envconfig.Process(context.Background(), &ttlDynmo); err != nil {
		log.Fatal(err.Error())
	}

	router := gin.Default()

	//router.Use(cors.Default())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "Cookie"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           1 * time.Hour,
	}))
	router.Use(gin.Logger())
	router.HandleMethodNotAllowed = true
	v1 := router.Group("/v1")
	{
		check := v1.Group("/check")
		{
			check.Use(middlewares.Auth)
			check.GET("/locate", cekdata.GetKec(db))
			check.GET("/match/:filename", cekdata.CheckDPT(db))
		}

		file := v1.Group("/files")
		{
			file.Use(middlewares.Auth)
			file.POST("/import", files.ImportExcel(db))
			file.GET("/export/:filename", files.ExportMatchExcel(db))
			file.GET("/read/:filename", files.ReadFile(db))
			file.GET("/list", files.GetListFolder)
			file.GET("/download/:filename", files.DownloadSampleXlsx)
			file.DELETE("/:filename", files.DeleteFile(db))
		}

		authentication := v1.Group("/auth")
		{
			authentication.POST("/register", auth.Register)
			authentication.POST("/login", auth.Login)
			authentication.GET("/logout", auth.Logout)
			authentication.POST("/forgot", auth.ForgotPassword)
			authentication.POST("/reset", auth.ResetPassword)
			authentication.POST("/confirm", auth.Confirmation)
			authentication.POST("/resend", auth.ResendCode)
		}
	}

	return router
}
