package cekdata

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"cekkustomer.com/api/models"
	"cekkustomer.com/dtos"
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

func CheckDPT(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var cekMatch models.ImportCustomerXls

		getKec, err := models.GetAllKec(db)
		if err != nil {
			log.Println(err.Error())
			return
		}

		// readfile excel
		fileName := ctx.Param("filename")
		folderUser := ctx.Param("foldername")
		filePath := fmt.Sprintf("%s/%s", folderUser, fileName)

		firstName := ctx.Query("nama")
		agenciesName := "folder-user"
		results := make(map[string]interface{})

		for _, tableName := range getKec {
			var err error
			var result []dtos.CheckDPT

			if firstName == "" {
				result, err = cekMatch.GetAll(db, tableName, agenciesName, filePath)
			} else {
				result, err = cekMatch.GetAllByName(db, tableName, agenciesName, firstName, filePath)
			}
			if err != nil {
				ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
				return
			}

			if len(result) > 0 {
				results[tableName] = result
			}
		}

		ctx.IndentedJSON(http.StatusOK, gin.H{
			"results": results,
		})

	}
}
