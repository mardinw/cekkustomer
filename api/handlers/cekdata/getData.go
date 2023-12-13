package cekdata

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"cekkustomer.com/api/middlewares"
	"cekkustomer.com/api/models"
	"cekkustomer.com/pkg/aws"
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

		getKec, err := models.GetAllKec(db)
		if err != nil {
			log.Println(err.Error())
			return
		}

		results := make(map[string]interface{})

		// readfile excel
		fileName := ctx.Param("filename")
		folderUser := ctx.Param("foldername")

		bucketName := "importxclxit"

		filePath := fmt.Sprintf("%s/%s", folderUser, fileName)

		getFile, err := aws.NewConnect().S3.GetFile(bucketName, filePath)
		if err != nil {
			log.Println(err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		readFile, err := middlewares.ReadExcel(getFile.Body)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// definisikan array

		//
		for _, tableName := range getKec {
			for _, data := range readFile {
				concatCustomerValue, ok := data["concat_customer"].(string)
				if !ok {
					log.Println("concat customer not found")
					continue
				}
				result, err := models.CheckData(db, tableName, concatCustomerValue)
				if err != nil {
					ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
					return
				}

				if len(result) > 0 {
					data["db_match"] = result
					results[tableName] = data
				}
			}
		}

		ctx.IndentedJSON(http.StatusOK, gin.H{
			"results": results,
		})
	}
}
