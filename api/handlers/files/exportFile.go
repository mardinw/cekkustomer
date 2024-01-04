package files

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"cekkustomer.com/api/middlewares"
	"cekkustomer.com/api/models"
	"cekkustomer.com/configs"
	"cekkustomer.com/pkg/aws"
	"github.com/gin-gonic/gin"
	"github.com/sethvargo/go-envconfig"
	"golang.org/x/net/context"
)

func ExportMatchExcel(db *sql.DB) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		var cekMatch models.ImportCustomerXls
		var bucketFolder configs.AwsS3Bucket
		if err := envconfig.Process(context.Background(), &bucketFolder); err != nil {
			log.Fatal(err.Error())
		}

		uuid, exists := ctx.Get("uuid")
		if !exists {
			log.Println("uuid tidak ditemukan")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		uuidStr, ok := uuid.(string)
		if !ok {
			log.Println("gagal konversi ke string")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		getKec, err := models.GetAllKec(db)
		if err != nil {
			log.Println(err.Error())
			return
		}

		fileName := ctx.Param("filename")

		agenciesName := uuidStr

		filePath := fmt.Sprintf("%s/%s", uuidStr, fileName)

		results := make(map[string]interface{})
		// get query di match
		for _, tableName := range getKec {
			result, err := cekMatch.GetMatchConcat(db, tableName, agenciesName, filePath)
			if err != nil {
				ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
				return
			}

			// convert the result to a map

			resultMap := make(map[string]interface{})
			resultMap[tableName] = result
			results[tableName] = resultMap
		}

		if len(results) > 0 {

			jsonData, err := json.Marshal(results)
			if err != nil {
				log.Println(err)
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			if err := middlewares.CreateExcel(string(jsonData), bucketFolder.ExportS3, fileName, filePath); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		if err := aws.NewConnect().S3.DownloadFile(bucketFolder.ExportS3, filePath, fileName); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer func() {
			err := os.Remove(fileName)
			if err != nil {
				log.Println(err.Error())
			}
		}()

		ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))

		ctx.File(fileName)
	}
}
