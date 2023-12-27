package files

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"cekkustomer.com/api/models"
	"cekkustomer.com/configs"
	"cekkustomer.com/pkg/aws"
	"github.com/gin-gonic/gin"
	"github.com/sethvargo/go-envconfig"
	"golang.org/x/net/context"
)

func DeleteFile(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var customer models.ImportCustomerXls
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

		fileName := ctx.Param("filename")
		agenciesName := uuidStr

		filePath := fmt.Sprintf("%s/%s", uuidStr, fileName)

		if err := aws.NewConnect().S3.DeleteFile(bucketFolder.ImportS3, filePath); err != nil {
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
