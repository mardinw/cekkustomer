package files

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"cekkustomer.com/api/models"
	"cekkustomer.com/configs"
	"github.com/gin-gonic/gin"
	"github.com/sethvargo/go-envconfig"
	"golang.org/x/net/context"
)

func ReadFile(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
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

		filePath := fmt.Sprintf("%s/%s", uuidStr, fileName)
		agenciesName := uuidStr

		var dataPreview models.ImportCustomerXls

		result, err := dataPreview.GetCustomer(db, filePath, agenciesName)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}

		// read file from bucket
		ctx.JSON(http.StatusOK, result)

	}

}
