package files

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"cekkustomer.com/api/middlewares"
	"cekkustomer.com/api/models"
	"cekkustomer.com/configs"
	"cekkustomer.com/pkg/aws"
	"github.com/gin-gonic/gin"
	"github.com/sethvargo/go-envconfig"
	"golang.org/x/net/context"
)

func isAllowedExtension(fileName string, allowedExtension []string) bool {
	for _, ext := range allowedExtension {
		if strings.HasSuffix(fileName, ext) {
			return true
		}
	}

	return false
}

func ImportExcel(db *sql.DB) gin.HandlerFunc {
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

		localUploadDir := "./uploads"
		agenciesName := uuidStr

		file, err := ctx.FormFile("file")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fileName := file.Filename

		// check file
		allowedExtension := []string{".xlsx", ".xls"}
		if !isAllowedExtension(fileName, allowedExtension) {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "file format only xlsx and xls"})
			return
		}

		filePath := filepath.Join(localUploadDir, fileName)

		if err := ctx.SaveUploadedFile(file, filePath); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		s3FilePath := filepath.Join(uuidStr, fileName)

		// upload file to s3
		if err := aws.NewConnect().S3.UploadFile(bucketFolder.ImportS3, filePath, s3FilePath); err != nil {
			log.Println("Failed to upload file to S3:", err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := os.Remove(filePath); err != nil {
			log.Println("Failed to remove uploaded file:", err.Error())
		} else {
			log.Println("File removed successfully:", filePath)
		}

		// read file from bucket
		getFile, err := aws.NewConnect().S3.GetFile(bucketFolder.ImportS3, s3FilePath)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		readFile, err := middlewares.ReadExcel(getFile.Body, bucketFolder.ImportS3, s3FilePath)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// end readfile
		// start insert to table customer

		for _, data := range readFile {
			timeNow := time.Now().UnixMilli()

			cardNumber, err := strconv.ParseInt(data["card_number"].(string), 10, 64)
			if err != nil {
				log.Println(err)
				continue
			}

			concatCustomerValue, ok := data["concat_customer"].(string)
			if !ok {
				log.Println("concat customer not found")
				return
			}

			concatCustToUpper := strings.ToUpper(concatCustomerValue)

			inputCustomer := &models.ImportCustomerXls{
				CardNumber:     cardNumber,
				FirstName:      data["first_name"].(string),
				Collector:      data["collector"].(string),
				Agencies:       agenciesName,
				Address3:       data["address_3"].(string),
				Address4:       data["address_4"].(string),
				ZipCode:        data["zipcode"].(string),
				ConcatCustomer: concatCustToUpper,
				Files:          s3FilePath,
				Created:        timeNow,
			}

			if err := inputCustomer.InsertCustomer(db); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
				return
			}

			ctx.JSON(http.StatusOK, gin.H{
				"message": "successfully uploaded",
			})
		}
	}
}
