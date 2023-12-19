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
	"cekkustomer.com/pkg/aws"
	"github.com/gin-gonic/gin"
)

func ExportMatchExcel(db *sql.DB) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		var cekMatch models.ImportCustomerXls

		getKec, err := models.GetAllKec(db)
		if err != nil {
			log.Println(err.Error())
			return
		}

		fileName := ctx.Param("filename")
		folderUser := ctx.Param("foldername")

		agenciesName := "folder-user"

		bucketExport := "exportxclxit"

		filePath := fmt.Sprintf("%s/%s", folderUser, fileName)

		results := make(map[string]interface{})
		// get query di match
		for _, tableName := range getKec {
			result, err := cekMatch.GetAll(db, tableName, agenciesName, filePath)
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

			if err := middlewares.CreateExcel(string(jsonData), bucketExport, fileName, filePath); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		if err := aws.NewConnect().S3.DownloadFile(bucketExport, filePath, fileName); err != nil {
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
