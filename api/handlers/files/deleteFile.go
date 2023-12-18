package files

import (
	"database/sql"
	"fmt"
	"net/http"

	"cekkustomer.com/api/models"
	"cekkustomer.com/pkg/aws"
	"github.com/gin-gonic/gin"
)

func DeleteFile(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var customer models.ImportCustomerXls

		bucketName := "importxclxit"
		fileName := ctx.Param("filename")
		folderUser := ctx.Param("foldername")
		agenciesName := "folder-user"

		filePath := fmt.Sprintf("%s/%s", folderUser, fileName)

		if err := aws.NewConnect().S3.DeleteFile(bucketName, filePath); err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
			return
		}

		if err := customer.DeleteCustomer(db, filePath, agenciesName); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("delete file %s successfully", filePath)})
	}
}
