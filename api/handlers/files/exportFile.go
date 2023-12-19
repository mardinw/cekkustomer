package files

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"cekkustomer.com/api/middlewares"
	"cekkustomer.com/api/models"
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
		//bucketName := "importxclxit"

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

			if len(result) > 0 {
				results[tableName] = result
			}
		}

		if len(results) > 0 {
			jsonData, err := json.Marshal(results)
			if err != nil {
				log.Println(err)
				return
			}

			if err := middlewares.CreateExcel(string(jsonData), bucketExport, filePath); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "file successfully create"})
	}
}
