package files

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"cekkustomer.com/api/models"
	"cekkustomer.com/configs"
	"cekkustomer.com/dtos"
	"cekkustomer.com/pkg/aws"
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
		var result []dtos.DataPreviewNIK
		var err error

		if !aws.NewConnect().S3.CheckExists(ctx, bucketFolder.ImportS3, filePath) {
			ctx.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("file %s not found", filePath)})
			return
		}
		if firstName == "" {
			result, _ = dataPreview.GetCustomer(db, filePath, agenciesName)
		} else {
			result, err = dataPreview.GetCustomerByName(db, agenciesName, firstName, filePath)
			if err != nil {
				ctx.JSON(http.StatusNotFound, gin.H{"message": "data not found"})
				return
			}
		}

		if len(result) == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "data not found"})
			return
		}

		// var nikList []string
		// for _, preview := range result {
		// 	nikList = append(nikList, preview.NIK)
		// }

		// nikLengths := make(map[string]int)
		// for _, nik := range nikList {
		// 	nikLengths[nik] = len(nik)
		// }

		// read file from bucket
		ctx.JSON(http.StatusOK, result)

	}

}
