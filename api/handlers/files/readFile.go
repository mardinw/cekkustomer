package files

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"cekkustomer.com/api/models"
	"cekkustomer.com/configs"
	"cekkustomer.com/dtos"
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

		// search nama di table customer
		firstName := ctx.Query("first_name")

		var dataPreview models.ImportCustomerXls
		var result []dtos.DataPreview
		var err error

		if firstName == "" {
			result, err = dataPreview.GetCustomer(db, filePath, agenciesName)
		} else {
			result, err = dataPreview.GetCustomerByName(db, filePath, agenciesName, firstName)
		}
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}

		// read file from bucket
		ctx.JSON(http.StatusOK, result)

	}

}
