package cekdata

import (
	"database/sql"
	"net/http"

	"cekkustomer.com/api/models"
	"github.com/gin-gonic/gin"
)

func GetKec(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		result, err := models.GetAllKec(db)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"locate": result})
	}
}

func GetDPT(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var dpt models.DPT

		//getKec, err := models.GetAllKec(db)
		//if err != nil {
		//	log.Println(err.Error())
		//	return
		//}

		results, err := dpt.GetAll(db, "dpt_kiaracondong")
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"results": results})
	}
}
