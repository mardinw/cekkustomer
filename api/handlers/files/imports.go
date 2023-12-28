package files

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
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

		dataChannel := make(chan map[string]interface{}, 10)

		readFile, err := middlewares.ReadExcel(getFile.Body, bucketFolder.ImportS3, s3FilePath)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// end readfile
		// start insert to table customer

		// goroutine for import data concurrently
		go func() {
			for _, data := range readFile {
				dataChannel <- data
			}

			close(dataChannel)
		}()

		// Goroutines to process data concurrently
		var wg sync.WaitGroup
		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				processData(db, agenciesName, s3FilePath, dataChannel)
			}()
		}

		wg.Wait()

		ctx.JSON(http.StatusOK, gin.H{
			"message": "successfully uploaded",
		})
	}
}

func processData(db *sql.DB, agenciesName, s3FilePath string, dataChannel <-chan map[string]interface{}) {
	for data := range dataChannel {
		timeNow := time.Now().UnixMilli()

		// card number
		cardNumberRaw, ok := data["card_number"]
		var cardNumber int64
		var err error

		if !ok || cardNumberRaw == nil {
			log.Println("card number not found or nil")
			cardNumber = int64(0)
		} else {
			cardNumberStr, ok := cardNumberRaw.(string)
			if !ok {
				log.Println("card_number is not a string")
				continue
			}
			cardNumber, err = strconv.ParseInt(cardNumberStr, 10, 64)
			if err != nil {
				log.Println("card_number not found or is nil")
				cardNumber = int64(0)
			}
		}

		// collector
		collectorRaw, ok := data["collector"]
		var collector string

		if !ok || collectorRaw == nil {
			log.Println("collector not found or nil")
			collector = ""
		} else {
			collector, ok = collectorRaw.(string)
			if !ok {
				log.Println("collector is not a string")
				continue
			}
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
			Collector:      collector,
			Agencies:       agenciesName,
			Address3:       data["address_3"].(string),
			Address4:       data["address_4"].(string),
			ZipCode:        data["zipcode"].(string),
			ConcatCustomer: concatCustToUpper,
			Files:          s3FilePath,
			Created:        timeNow,
		}

		if err := inputCustomer.InsertCustomer(db); err != nil {
			log.Println("failed to insert customer:", err.Error())
			return
		}

	}
}
