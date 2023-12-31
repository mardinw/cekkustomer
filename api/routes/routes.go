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
			check.GET("/match/concat/:filename", cekdata.CheckDPTByConcat(db))
			check.GET("/match/nik/:filename", cekdata.CheckDPTByNIK(db))
			check.GET("/attribute/:user", cekdata.GetAttributes)
		}

		file := v1.Group("/files")
		file.GET("/download/:filename", files.DownloadSampleXlsx)
		{
			file.Use(middlewares.Auth)
			file.POST("/import", files.ImportExcel(db))
			file.GET("/export/concat/:filename", files.ExportMatchConcatExcel(db))
			file.GET("/export/nik/:filename", files.ExportMatchNIKExcel(db))
			file.GET("/read/:filename", files.ReadFile(db))
			file.GET("/list", files.GetListFolder)
			file.DELETE("/:filename", files.DeleteFile(db))
		}

		authentication := v1.Group("/auth")
		{
			authentication.GET("/check", middlewares.Auth)
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
