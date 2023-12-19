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
	"cekkustomer.com/pkg/aws"
	"github.com/gin-gonic/gin"
)

func ImportExcel(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		localUploadDir := "./uploads"
		uploadFolder := "folder-user"
		bucketName := "importxclxit"
		agenciesName := "folder-user"

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

		s3FilePath := filepath.Join(uploadFolder, fileName)

		// upload file to s3
		if err := aws.NewConnect().S3.UploadFile(bucketName, filePath, s3FilePath); err != nil {
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
		getFile, err := aws.NewConnect().S3.GetFile(bucketName, s3FilePath)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		readFile, err := middlewares.ReadExcel(getFile.Body)
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

			concatCustomerValue, ok := data["concat_customer (nama + tgl lahir)"].(string)
			if !ok {
				log.Println("concat customer not found")
				continue
			}
			concatCustToUpper := strings.ToUpper(concatCustomerValue)

			inputCustomer := &models.ImportCustomerXls{
				CardNumber:     cardNumber,
				FirstName:      data["first_name"].(string),
				Collector:      data["collector"].(string),
				Agencies:       agenciesName,
				Address3:       data["address_3"].(string),
				Address4:       data["address_4"].(string),
				ZipCode:        data["home_zip_code"].(string),
				ConcatCustomer: concatCustToUpper,
				Files:          s3FilePath,
				Created:        timeNow,
			}
			if err := inputCustomer.InsertCustomer(db); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
				return
			}
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "successfully insert",
		})

	}
}

func isAllowedExtension(fileName string, allowedExtension []string) bool {
	for _, ext := range allowedExtension {
		if strings.HasSuffix(fileName, ext) {
			return true
		}
	}

	return false
}
