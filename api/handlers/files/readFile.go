package files

import (
	"database/sql"
	"fmt"
	"net/http"

	"cekkustomer.com/api/models"
	"github.com/gin-gonic/gin"
)

func ReadFile(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//bucketName := "importxclxit"
		fileName := ctx.Param("filename")
		folderUser := ctx.Param("foldername")

		filePath := fmt.Sprintf("%s/%s", folderUser, fileName)
		agenciesName := "folder-user"
		var dataPreview models.ImportCustomerXls

		result, err := dataPreview.GetCustomer(db, filePath, agenciesName)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, result)
	}
}
