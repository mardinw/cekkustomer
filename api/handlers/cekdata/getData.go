package cekdata

import (
	"database/sql"
	"net/http"

	"cekkustomer.com/api/models"
	"github.com/gin-gonic/gin"
)

func GetDPT(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var dpt models.DPT

		result, err := dpt.GetAll(db)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"dpt_kiaracondong": result})
	}
}
